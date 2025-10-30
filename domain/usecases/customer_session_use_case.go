package usecases

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/responses"
	"golang.org/x/crypto/bcrypt"
)

type CustomerSessionUseCase struct {
	sessionRepo repositories.CustomerSessionRepository
	userRepo    repositories.UserRepo
	db          *sql.DB
}

func ProvideCustomerSessionUseCase(
	sessionRepo repositories.CustomerSessionRepository,
	userRepo repositories.UserRepo,
	db *sql.DB,
) *CustomerSessionUseCase {
	return &CustomerSessionUseCase{
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
		db:          db,
	}
}

// RegisterCustomerSession - Use Case 2.2C: ลงทะเบียนผู้ใช้งานสำหรับคอร์ส Sessions
func (u *CustomerSessionUseCase) RegisterCustomerSession(ctx context.Context, req requests.RegisterCustomerSessionRequest) (*responses.RegisterCustomerSessionResponse, error) {
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

// CheckBookingPermission - ตรวจสอบสิทธิ์การเข้าถึงฟังก์ชันการจอง
// ตรวจว่า Customer มีแพ็กเกจ Sessions แบบ ACTIVE และยังมีสิทธิ์คงเหลือ
func (u *CustomerSessionUseCase) CheckBookingPermission(ctx context.Context, customerUsername string) (bool, error) {
	hasPermission, err := u.sessionRepo.CheckBookingPermission(ctx, customerUsername)
	if err != nil {
		return false, fmt.Errorf("failed to check booking permission: %w", err)
	}

	return hasPermission, nil
}
