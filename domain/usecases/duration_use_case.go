package usecases

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/responses"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
	"golang.org/x/crypto/bcrypt"
)

type CustomerDurationUseCase struct {
	durationRepo repositories.CustomerDurationRepository
	productRepo  repositories.ProductRepository
	userRepo     repositories.UserRepo
	db           *sql.DB
}

func ProvideCustomerDurationUseCase(
	durationRepo repositories.CustomerDurationRepository,
	productRepo repositories.ProductRepository,
	userRepo repositories.UserRepo,
	db *sql.DB,
) *CustomerDurationUseCase {
	return &CustomerDurationUseCase{
		durationRepo: durationRepo,
		productRepo:  productRepo,
		userRepo:     userRepo,
		db:           db,
	}
}

// PurchaseDuration creates a new customer duration purchase
// Flow: ตรง use case description step 4
func (u *CustomerDurationUseCase) PurchaseDuration(ctx context.Context, username string, req requests.PurchaseDurationRequest) (int32, error) {
	// Get product details to calculate end date
	product, err := u.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return 0, fmt.Errorf("product not found: %w", err)
	}

	// Validate product type
	if product.Type != "DURATION" {
		return 0, fmt.Errorf("product is not a DURATION type")
	}

	// Validate duration days
	if product.DurationDays == nil || *product.DurationDays <= 0 {
		return 0, fmt.Errorf("invalid product: duration_days not set")
	}

	// Calculate dates
	now := time.Now()
	startDate := now
	endDate := startDate.AddDate(0, 0, int(*product.DurationDays))

	// Format prices as DECIMAL strings
	pricePaidStr := fmt.Sprintf("%.2f", req.PricePaid)
	discountAmountStr := fmt.Sprintf("%.2f", req.DiscountAmount)

	// Create duration record
	params := repositories.CreateCustomerDurationParams{
		CustomerUsername: username,
		SalesUsername:    req.SalesUsername,
		ProductID:        req.ProductID,
		PurchaseDate:     now,
		StartDate:        startDate,
		EndDate:          endDate,
		PricePaid:        pricePaidStr,
		DiscountAmount:   discountAmountStr,
		Status:           "ACTIVE",
	}

	if err := u.durationRepo.Create(ctx, params); err != nil {
		return 0, fmt.Errorf("failed to create duration: %w", err)
	}

	// Return success (in real implementation, we'd return the created ID)
	return product.ID, nil
}

// GetCustomerDurations returns all durations for a customer
func (u *CustomerDurationUseCase) GetCustomerDurations(ctx context.Context, username string) ([]responses.CustomerDurationResponse, error) {
	durations, err := u.durationRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer durations: %w", err)
	}

	result := make([]responses.CustomerDurationResponse, len(durations))
	for i, d := range durations {
		result[i] = u.mapToResponse(d)
	}
	return result, nil
}

// Helper to convert repository model to response
func (u *CustomerDurationUseCase) mapToResponse(d repositories.CustomerDurationInfo) responses.CustomerDurationResponse {
	pricePaid, _ := strconv.ParseFloat(d.PricePaid, 64)
	discountAmount, _ := strconv.ParseFloat(d.DiscountAmount, 64)

	return responses.CustomerDurationResponse{
		ID:               d.ID,
		CustomerUsername: d.CustomerUsername,
		SalesUsername:    d.SalesUsername,
		ProductID:        d.ProductID,
		PurchaseDate:     d.PurchaseDate,
		StartDate:        d.StartDate,
		EndDate:          d.EndDate,
		PricePaid:        pricePaid,
		DiscountAmount:   discountAmount,
		Status:           d.Status,
		CreatedAt:        d.CreatedAt,
		UpdatedAt:        d.UpdatedAt,
	}
}

func (uc *CustomerDurationUseCase) List(ctx context.Context, req requests.ListCustomerDurationsRequest) (responses.ListCustomerDurationsResponse, error) {
	limit := req.Limit
	if limit <= 0 || limit > 100 { limit = 10 }
	page := req.Page
	if page <= 0 { page = 1 }
	offset := (page - 1) * limit

	data, err := uc.durationRepo.List(ctx, limit, offset)
	if err != nil { return responses.ListCustomerDurationsResponse{}, err }
	total, err := uc.durationRepo.Count(ctx)
	if err != nil { return responses.ListCustomerDurationsResponse{}, err }

	if len(data) == 0 {
		data = []dbmodel.ListCustomerDurationsRow{}
	}

	return responses.ListCustomerDurationsResponse{
		Data: data,
		Meta: responses.PageMeta{
			Page:       page,
			Limit:      limit,
			TotalItems: total,
			TotalPages: int32(math.Ceil(float64(total)/float64(limit))),
		},
	}, nil
}

var reYMD = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
func (uc *CustomerDurationUseCase) UpdateDuration(
	ctx context.Context,
	id int32,
	req requests.UpdateCustomerDurationRequest,
) (responses.CustomerDurationUpdatedResponse, error) {

	// 1) Validate input
	if !reYMD.MatchString(req.StartDate) {
		return responses.CustomerDurationUpdatedResponse{}, errors.New("invalid startDate (YYYY-MM-DD)")
	}
	if req.PricePaid < 0 {
		return responses.CustomerDurationUpdatedResponse{}, errors.New("pricePaid must be >= 0")
	}
	// DiscountAmount เป็น float64 (ไม่ใช่ *float64) => เช็คค่าตรง ๆ
	if req.DiscountAmount < 0 {
		return responses.CustomerDurationUpdatedResponse{}, errors.New("discountAmount must be >= 0")
	}
	switch req.Status {
	case "ACTIVE", "EXPIRED", "FROZEN", "CANCELLED":
	default:
		return responses.CustomerDurationUpdatedResponse{}, errors.New("invalid status")
	}

	// 2) ตรวจสอบว่า product ที่อ้างอิงเป็น DURATION และมี duration_days > 0
	days, err := uc.durationRepo.GetDurationDaysForDurationID(ctx, id)
	if err != nil {
		return responses.CustomerDurationUpdatedResponse{}, fmt.Errorf("duration not found or product invalid: %w", err)
	}
	if days <= 0 {
		return responses.CustomerDurationUpdatedResponse{}, errors.New("invalid product duration_days")
	}

	// 3) แปลงราคาเป็น string ตามที่ repo/db ต้องการ
	pricePaidStr := fmt.Sprintf("%.2f", req.PricePaid)
	discountStr := fmt.Sprintf("%.2f", req.DiscountAmount)

	// 4) Persist (SQL จะคำนวณ end_date = start + duration_days - 1)
	err = uc.durationRepo.UpdateEditableFields(ctx, repositories.UpdateCustomerDurationEditableFieldsParams{
		ID:             id,
		StartDateYMD:   req.StartDate,     // ← ใช้ D ใหญ่ ให้ตรง struct
		PricePaid:      pricePaidStr,
		DiscountAmount: &discountStr,      // ← domain param เป็น *string
		Status:         req.Status,
	})
	if err != nil {
		return responses.CustomerDurationUpdatedResponse{}, err
	}

	return responses.CustomerDurationUpdatedResponse{
		Message: fmt.Sprintf("Duration Package ID: %d updated successfully", id),
	}, nil
}

func (uc *CustomerDurationUseCase) Delete(
	ctx context.Context,
	id int32,
) (responses.CustomerDurationDeletedResponse, error) {

	if err := uc.durationRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return responses.CustomerDurationDeletedResponse{}, fmt.Errorf("duration package not found")
		}
		return responses.CustomerDurationDeletedResponse{}, err
	}

	return responses.CustomerDurationDeletedResponse{
		Message: fmt.Sprintf("Package ID: %d deleted successfully", id),
	}, nil
}
// RegisterCustomerDuration - Use Case 2.1C: ลงทะเบียนผู้ใช้งานสำหรับแพ็กเกจ Duration
func (u *CustomerDurationUseCase) RegisterCustomerDuration(ctx context.Context, req requests.RegisterCustomerDurationRequest) (*responses.RegisterCustomerDurationResponse, error) {
	// Validate username not exists (Q2.1C.1)
	exists, err := u.userRepo.CheckUsernameExists(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("USERNAME_ALREADY_EXISTS")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Parse date of birth
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return nil, fmt.Errorf("invalid date format for dateOfBirth, expected YYYY-MM-DD: %w", err)
	}

	// Parse start date
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date format for startDate, expected YYYY-MM-DD: %w", err)
	}

	// Calculate end date (StartDate + DurationDays)
	endDate := startDate.AddDate(0, 0, int(req.DurationDays))

	// Convert prices to DECIMAL strings
	pricePaidStr := fmt.Sprintf("%.2f", req.PricePaid)
	discountAmountStr := fmt.Sprintf("%.2f", req.DiscountAmount)

	// Prepare parameters
	params := repositories.RegisterDurationParams{
		// User
		Username:    req.Username,
		Password:    string(hashedPassword),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Gender:      req.Gender,
		DateOfBirth: dob,
		PhoneNumber: req.PhoneNumber,
		Gmail:       strings.ToLower(req.Gmail), // Lowercase as per spec
		// Customer
		HealthInfo:                   req.HealthInfo,
		Address:                      req.Address,
		CompanyName:                  req.CompanyName,
		CompanyPosition:              req.CompanyPosition,
		MaritalStatus:                req.MaritalStatus,
		EmergencyContactName:         req.EmergencyContactName,
		EmergencyContactRelationship: req.EmergencyContactRelationship,
		EmergencyContactPhone:        req.EmergencyContactPhone,
		MarketingSource:              req.MarketingSource,
		// CustomerDuration
		ProductID:      req.ProductID,
		SalesUsername:  req.SalesUsername,
		PurchaseDate:   time.Now(),
		StartDate:      startDate,
		EndDate:        endDate,
		PricePaid:      pricePaidStr,
		DiscountAmount: discountAmountStr,
	}

	// Begin transaction
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Execute registration (creates User, Customer, CustomerDuration)
	durationID, err := u.durationRepo.RegisterDuration(ctx, tx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to register customer duration: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Build response
	return &responses.RegisterCustomerDurationResponse{
		Username:       req.Username,
		DurationID:     durationID,
		ProductID:      req.ProductID,
		SalesUsername:  req.SalesUsername,
		StartDate:      startDate,
		EndDate:        endDate,
		DurationDays:   req.DurationDays,
		PricePaid:      pricePaidStr,
		DiscountAmount: discountAmountStr,
		Message:        "Customer duration registered successfully",
	}, nil
}
