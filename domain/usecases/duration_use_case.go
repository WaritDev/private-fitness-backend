package usecases

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/responses"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
	"github.com/WaritDev/private-fitness-backend/utils"
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

// GetCustomerActiveDuration - ดึงข้อมูล Duration packages ที่ ACTIVE ของลูกค้า (คล้าย GetCustomerActiveSessions)
func (u *CustomerDurationUseCase) GetCustomerActiveDuration(ctx context.Context, username string) ([]responses.CustomerDurationPackageResponse, error) {
	// เรียก repository เพื่อดึงข้อมูล
	packages, err := u.durationRepo.GetCustomerActiveDuration(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get active duration packages: %w", err)
	}

	// แปลงเป็น response DTO
	result := make([]responses.CustomerDurationPackageResponse, 0, len(packages))
	for _, pkg := range packages {
		// Convert DECIMAL strings to float64
		pricePaid, _ := strconv.ParseFloat(pkg.PricePaid, 64)
		discountAmount, _ := strconv.ParseFloat(pkg.DiscountAmount, 64)

		result = append(result, responses.CustomerDurationPackageResponse{
			ID:               pkg.ID,
			CustomerUsername: pkg.CustomerUsername,
			ProductID:        pkg.ProductID,
			ProductName:      pkg.ProductName,
			DurationDays:     pkg.DurationDays,
			SalesUsername:    pkg.SalesUsername,
			PurchaseDate:     pkg.PurchaseDate,
			StartDate:        pkg.StartDate,
			EndDate:          pkg.EndDate,
			DaysRemaining:    pkg.DaysRemaining,
			PricePaid:        pricePaid,
			DiscountAmount:   discountAmount,
			Status:           pkg.Status,
			CreatedAt:        pkg.CreatedAt,
		})
	}

	return result, nil
}

func (uc *CustomerDurationUseCase) List(ctx context.Context) ([]dbmodel.ListCustomerDurationsRow, error) {
	return uc.durationRepo.List(ctx)
}

func (uc *CustomerDurationUseCase) UpdateDuration(
	ctx context.Context,
	id int32,
	req requests.UpdateCustomerDurationRequest,
) (responses.CustomerDurationUpdatedResponse, error) {
	switch strings.ToUpper(req.Status) {
	case "ACTIVE", "EXPIRED", "FROZEN", "CANCELLED":
	default:
		return responses.CustomerDurationUpdatedResponse{}, errors.New("invalid status")
	}

	days, err := uc.durationRepo.GetDurationDaysForDurationID(ctx, id)
	if err != nil {
		return responses.CustomerDurationUpdatedResponse{}, fmt.Errorf("duration not found or product invalid: %w", err)
	}
	if days <= 0 {
		return responses.CustomerDurationUpdatedResponse{}, errors.New("invalid product duration_days")
	}

	pricePaid := round2(req.PricePaid)
	discount := round2(req.DiscountAmount)

	err = uc.durationRepo.UpdateEditableFields(ctx, repositories.UpdateCustomerDurationEditableFieldsParams{
		ID:             id,
		StartDateYMD:   req.StartDate,
		PricePaid:      pricePaid,
		DiscountAmount: &discount,
		Status:         req.Status,
	})
	if err != nil {
		return responses.CustomerDurationUpdatedResponse{}, err
	}

	return responses.CustomerDurationUpdatedResponse{
		Message: fmt.Sprintf("Duration Package ID: %d updated successfully", id),
	}, nil
}

func round2(f float64) float64 {
	return math.Round(f*100) / 100
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

func (uc *CustomerDurationUseCase) GetByID(ctx context.Context, id string) (responses.CustomerDuration, error) {
	if strings.TrimSpace(id) == "" {
		return responses.CustomerDuration{}, errors.New("id required")
	}

	id32, err := utils.Atoi32(id)
	if err != nil {
		return responses.CustomerDuration{}, errors.New("invalid id")
	}

	row, err := uc.durationRepo.GetByID(ctx, id32)
	if err != nil {
		return responses.CustomerDuration{}, err
	}

	pricePaid := parseToFloat2(row.PricePaid)
	discount := parseToFloat2(row.DiscountAmount)

	resp := responses.CustomerDuration{
		ID:               utils.Itoa(row.ID),
		CustomerUsername: row.CustomerUsername,
		ProductID:        utils.Itoa(row.ProductID),
		SalesUsername:    utils.PtrToString(row.SalesUsername),
		PurchaseDate:     utils.ToYMD(row.PurchaseDate),
		StartDate:        utils.ToYMD(row.StartDate),
		EndDate:          utils.ToYMD(row.EndDate),
		PricePaid:        pricePaid,
		DiscountAmount:   discount,
		Status:           strings.ToUpper(string(row.Status)),
	}
	return resp, nil
}

func parseToFloat2(s string) float64 {
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0
	}
	return math.Round(f*100) / 100
}

// RenewDuration - Customer self-purchase duration package (ลูกค้าซื้อเพิ่มเอง)
func (u *CustomerDurationUseCase) RenewDuration(ctx context.Context, customerUsername string, req requests.RenewDurationRequest) (*responses.RenewDurationResponse, error) {
	// Step 1: Validate product exists and is DURATION type
	product, err := u.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("PRODUCT_NOT_FOUND_OR_INACTIVE")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	// Step 2: Validate product type
	if product.Type != "DURATION" {
		return nil, fmt.Errorf("INVALID_PRODUCT_DURATION_AMOUNT")
	}

	// Step 3: Validate duration_days
	if product.DurationDays == nil || *product.DurationDays <= 0 {
		return nil, fmt.Errorf("INVALID_PRODUCT_DURATION_AMOUNT")
	}

	// Step 4: Parse list_price (DECIMAL string to float64)
	listPrice, err := strconv.ParseFloat(product.ListPrice, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse list price: %w", err)
	}

	// Step 5: Calculate start_date and end_date (SQL auto-calculates with NOW() + duration_days)
	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, int(*product.DurationDays))

	// Step 6: Calculate days_remaining (from today)
	daysRemaining := int32(endDate.Sub(time.Now()).Hours() / 24)
	if daysRemaining < 0 {
		daysRemaining = 0
	}

	// Step 7: INSERT new duration package
	// SQL will auto-calculate start_date = NOW() and end_date = NOW() + duration_days
	params := repositories.RenewDurationParams{
		CustomerUsername: customerUsername,
		ProductID:        req.ProductID,
		PricePaid:        product.ListPrice, // DECIMAL string
	}

	err = u.durationRepo.RenewDuration(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to renew duration: %w", err)
	}

	// Step 8: Return success response
	return &responses.RenewDurationResponse{
		ID:               0, // Will be updated after LAST_INSERT_ID() if needed
		CustomerUsername: customerUsername,
		ProductID:        req.ProductID,
		ProductName:      product.Name,
		DurationDays:     *product.DurationDays,
		SalesUsername:    nil, // NULL for self-purchase
		PurchaseDate:     time.Now(),
		StartDate:        startDate,
		EndDate:          endDate,
		DaysRemaining:    daysRemaining,
		PricePaid:        listPrice,
		DiscountAmount:   0, // No discount for self-purchase
		Status:           "ACTIVE",
		Message:          "Duration package renewed successfully",
	}, nil
}
