package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/responses"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
	"github.com/WaritDev/private-fitness-backend/utils"
	"golang.org/x/crypto/bcrypt"
)

type CustomerSessionService struct {
	sessionRepo repositories.CustomerSessionRepository
	userRepo    repositories.UserRepo
	db          *sql.DB
}

func ProvideCustomerSessionService(
	sessionRepo repositories.CustomerSessionRepository,
	userRepo repositories.UserRepo,
	db *sql.DB,
) *CustomerSessionService {
	return &CustomerSessionService{
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
		db:          db,
	}
}

// RegisterCustomerSession - Use Case 2.2C: ลงทะเบียนผู้ใช้งานสำหรับคอร์ส Sessions
func (u *CustomerSessionService) RegisterCustomerSession(ctx context.Context, req requests.RegisterCustomerSessionRequest) (*responses.RegisterCustomerSessionResponse, error) {
	// Validate username not exists (Q2.2C.1)
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
		return nil, fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}

	// Convert prices to DECIMAL strings
	pricePaidStr := fmt.Sprintf("%.2f", req.PricePaid)
	discountAmountStr := fmt.Sprintf("%.2f", req.DiscountAmount)

	// Calculate training schedules (each session = 2 hours)
	schedules := make([]repositories.TrainingScheduleParams, len(req.SessionSchedules))
	for i, s := range req.SessionSchedules {
		schedules[i] = repositories.TrainingScheduleParams{
			TrainerUsername:  req.TrainerUsername,
			CustomerUsername: req.Username,
			// SessionID will be set after creating customer_session
			StartTime:    s.StartTime,
			EndTime:      s.StartTime.Add(2 * time.Hour), // End = Start + 2 hours
			ScheduleType: "APPOINTMENT",
		}
	}

	// Prepare parameters
	params := repositories.RegisterSessionParams{
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
		// CustomerSession
		TrainerUsername: req.TrainerUsername,
		ProductID:       req.ProductID,
		SalesUsername:   req.SalesUsername,
		PurchaseDate:    time.Now(),
		TotalSessions:   int32(len(schedules)), // Total = number of schedules
		UsedSessions:    0,                     // Start with 0 used sessions
		PricePaid:       pricePaidStr,
		DiscountAmount:  discountAmountStr,
		// Training Schedules
		Schedules: schedules,
	}

	// Begin transaction
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Execute registration (creates User, Customer, CustomerSession, TrainingSchedules, CustomerLog)
	sessionID, err := u.sessionRepo.RegisterCustomerSession(ctx, tx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to register customer session: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Build response
	createdSchedules := make([]responses.TrainingScheduleInfo, len(schedules))
	for i, s := range schedules {
		dayOfWeek := req.SessionSchedules[i].DayOfWeek
		createdSchedules[i] = responses.TrainingScheduleInfo{
			ScheduleID: int32(i + 1), // Placeholder, actual ID would need SELECT
			StartTime:  s.StartTime,
			EndTime:    s.EndTime,
			DayOfWeek:  dayOfWeek,
		}
	}

	return &responses.RegisterCustomerSessionResponse{
		Username:         req.Username,
		SessionID:        sessionID,
		TrainerUsername:  req.TrainerUsername,
		ProductID:        req.ProductID,
		TotalSessions:    len(schedules),
		SchedulesCreated: len(schedules),
		CreatedSchedules: createdSchedules,
		Message:          "Customer session registered successfully",
	}, nil
}

// CheckBookingPermission - Q2C.1: ตรวจสอบสิทธิ์การเข้าถึงฟังก์ชันการจองก่อนโหลดปฏิทิน
// ตรวจสอบว่า Customer มีแพ็กเกจ Sessions แบบ ACTIVE หรือไม่
// หมายเหตุ: ถ้าทำครบแล้วจะเปลี่ยน status เป็น 'COMPLETED' โดยอัตโนมัติ
func (u *CustomerSessionService) CheckBookingPermission(ctx context.Context, customerUsername string) (bool, error) {
	hasPermission, err := u.sessionRepo.CheckBookingPermission(ctx, customerUsername)
	if err != nil {
		return false, fmt.Errorf("failed to check booking permission: %w", err)
	}

	return hasPermission, nil
}

// GetCustomerActiveSessions - ดึงข้อมูล Session packages ที่ ACTIVE ของลูกค้า
func (u *CustomerSessionService) GetCustomerActiveSessions(ctx context.Context, customerUsername string) ([]responses.CustomerSessionPackageResponse, error) {
	packages, err := u.sessionRepo.GetCustomerActiveSessions(ctx, customerUsername)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer active sessions: %w", err)
	}

	result := make([]responses.CustomerSessionPackageResponse, len(packages))
	for i, pkg := range packages {
		// Parse DECIMAL strings to float64
		var pricePaid, discountAmount float64
		fmt.Sscanf(pkg.PricePaid, "%f", &pricePaid)
		fmt.Sscanf(pkg.DiscountAmount, "%f", &discountAmount)

		result[i] = responses.CustomerSessionPackageResponse{
			ID:                pkg.ID,
			CustomerUsername:  pkg.CustomerUsername,
			TrainerUsername:   pkg.TrainerUsername,
			ProductID:         pkg.ProductID,
			ProductName:       pkg.ProductName,
			TotalSessions:     pkg.TotalSessions,
			UsedSessions:      pkg.UsedSessions,
			SessionsRemaining: pkg.SessionsRemaining,
			PurchaseDate:      pkg.PurchaseDate,
			PricePaid:         pricePaid,
			DiscountAmount:    discountAmount,
			Status:            pkg.Status,
			CreatedAt:         pkg.CreatedAt,
		}
	}

	return result, nil
}

func (uc *CustomerSessionService) List(ctx context.Context, req requests.ListCustomerSessionsRequest) (responses.ListCustomerSessionsResponse, error) {
	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	page := req.Page
	if page <= 0 {
		page = 1
	}
	rows, err := uc.sessionRepo.List(ctx)
	if err != nil {
		return responses.ListCustomerSessionsResponse{}, err
	}

	if rows == nil {
		rows = []dbmodel.ListCustomerSessionsRow{}
	}

	// Calculate total from rows length since List() returns all records
	total := int64(len(rows))

	return responses.ListCustomerSessionsResponse{
		Data: rows,
		Meta: responses.PageMeta{
			Page:       page,
			Limit:      limit,
			TotalItems: total,
			TotalPages: int32(math.Ceil(float64(total) / float64(limit))),
		},
	}, nil
}

func (uc *CustomerSessionService) Update(
	ctx context.Context,
	id int32,
	req requests.UpdateCustomerSessionRequest,
) (responses.CustomerSessionUpdatedResponse, error) {
	n, err := uc.sessionRepo.CheckTrainerExists(ctx, req.TrainerUsername)
	if err != nil {
		return responses.CustomerSessionUpdatedResponse{}, err
	}
	if n == 0 {
		return responses.CustomerSessionUpdatedResponse{}, errors.New("trainer not found or not a TRAINER")
	}

	params := repositories.UpdateCustomerSessionEditableFieldsParams{
		ID:              id,
		TrainerUsername: req.TrainerUsername,
		PricePaid:       sql.NullFloat64{Float64: req.PricePaid, Valid: true},
		DiscountAmount:  sql.NullFloat64{Float64: req.DiscountAmount, Valid: true},
		Status:          req.Status,
	}
	if err := uc.sessionRepo.UpdateEditableFields(ctx, params); err != nil {
		return responses.CustomerSessionUpdatedResponse{}, err
	}

	return responses.CustomerSessionUpdatedResponse{
		Message: fmt.Sprintf("Sessions Course ID: %d updated successfully", id),
	}, nil
}

func (uc *CustomerSessionService) Delete(
	ctx context.Context,
	id int32,
) (responses.CustomerSessionDeletedResponse, error) {

	if err := uc.sessionRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return responses.CustomerSessionDeletedResponse{}, fmt.Errorf("session course not found")
		}
		return responses.CustomerSessionDeletedResponse{}, err
	}

	return responses.CustomerSessionDeletedResponse{
		Message: fmt.Sprintf("Sessions Course ID: %d deleted successfully", id),
	}, nil
}

func (uc *CustomerSessionService) GetByID(ctx context.Context, id string) (responses.CustomerSession, error) {
	if strings.TrimSpace(id) == "" {
		return responses.CustomerSession{}, errors.New("id required")
	}

	id32, err := utils.Atoi32(id)
	if err != nil {
		return responses.CustomerSession{}, errors.New("invalid id")
	}

	row, err := uc.sessionRepo.GetByID(ctx, id32)
	if err != nil {
		return responses.CustomerSession{}, err
	}

	var pricePaid float64
	if row.PricePaid != "" {
		fmt.Sscanf(row.PricePaid, "%f", &pricePaid)
	}

	var discountAmount float64
	if row.DiscountAmount.Valid && row.DiscountAmount.String != "" {
		fmt.Sscanf(row.DiscountAmount.String, "%f", &discountAmount)
	}

	resp := responses.CustomerSession{
		ID:               utils.Itoa(row.ID),
		CustomerUsername: utils.NS(row.CustomerUsername),
		TrainerUsername:  utils.NS(row.TrainerUsername),
		SalesUsername:    utils.NS(row.SalesUsername),
		ProductID:        utils.Itoa(utils.NI32(row.ProductID)),
		PurchaseDate:     utils.ToYMD(row.PurchaseDate),
		TotalSessions:    row.TotalSessions,
		UsedSessions:     utils.NI32(row.UsedSessions),
		PricePaid:        pricePaid,
		DiscountAmount:   discountAmount,
		Status:           strings.ToUpper(string(row.Status)),
	}

	return resp, nil
}

// RenewSession - Use Case: ต่ออายุ/ซื้อเพิ่ม Session Package (ลูกค้าซื้อเอง)
func (u *CustomerSessionService) RenewSession(ctx context.Context, customerUsername string, req requests.RenewSessionRequest) (*responses.RenewSessionResponse, error) {
	// 1. ตรวจสอบว่า product มีอยู่และเป็น SESSION type
	var (
		productID     int32
		productName   string
		productType   string
		sessionAmount sql.NullInt32
		listPrice     string
		isActive      int8
	)

	err := u.db.QueryRowContext(ctx, `
		SELECT id, name, type, session_amount, list_price, is_active
		FROM products
		WHERE id = ? AND type = 'SESSION' AND is_active = 1
	`, req.ProductID).Scan(&productID, &productName, &productType, &sessionAmount, &listPrice, &isActive)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("PRODUCT_NOT_FOUND_OR_INACTIVE")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if !sessionAmount.Valid || sessionAmount.Int32 <= 0 {
		return nil, fmt.Errorf("INVALID_PRODUCT_SESSION_AMOUNT")
	}

	// 2. ตรวจสอบว่า trainer มีอยู่จริง
	var trainerExists int64
	err = u.db.QueryRowContext(ctx, `
		SELECT COUNT(username)
		FROM users
		WHERE username = ? AND role = 'TRAINER'
	`, req.TrainerUsername).Scan(&trainerExists)

	if err != nil {
		return nil, fmt.Errorf("failed to check trainer: %w", err)
	}

	if trainerExists == 0 {
		return nil, fmt.Errorf("TRAINER_NOT_FOUND")
	}

	// 3. Parse price_paid จาก product list_price
	var pricePaid float64
	fmt.Sscanf(listPrice, "%f", &pricePaid)

	// 4. INSERT session package ใหม่
	params := repositories.RenewSessionParams{
		CustomerUsername: customerUsername,
		TrainerUsername:  req.TrainerUsername,
		ProductID:        req.ProductID,
		TotalSessions:    sessionAmount.Int32,
		PricePaid:        fmt.Sprintf("%.2f", pricePaid),
	}

	err = u.sessionRepo.RenewSession(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to renew session: %w", err)
	}

	// 5. Return response
	return &responses.RenewSessionResponse{
		CustomerUsername: customerUsername,
		TrainerUsername:  req.TrainerUsername,
		ProductID:        req.ProductID,
		ProductName:      productName,
		TotalSessions:    sessionAmount.Int32,
		UsedSessions:     0,
		PurchaseDate:     time.Now(),
		PricePaid:        pricePaid,
		DiscountAmount:   0,
		Status:           "ACTIVE",
		Message:          "Session package renewed successfully",
	}, nil
}
