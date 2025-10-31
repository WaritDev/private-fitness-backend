package usecases

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/responses"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
	"github.com/WaritDev/private-fitness-backend/utils"
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

// GetCustomerActiveSessions - ดึงข้อมูล Session packages ที่ ACTIVE ของลูกค้า
func (u *CustomerSessionUseCase) GetCustomerActiveSessions(ctx context.Context, customerUsername string) ([]responses.CustomerSessionPackageResponse, error) {
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

func (uc *CustomerSessionUseCase) List(ctx context.Context) (responses.ListCustomerSessionsResponse, error) {
	rows, err := uc.sessionRepo.List(ctx)
	if err != nil {
		return responses.ListCustomerSessionsResponse{}, err
	}
	if rows == nil {
		rows = []dbmodel.ListCustomerSessionsRow{}
	}
	return responses.ListCustomerSessionsResponse{Data: rows}, nil
}


var reUsername = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]{2,29}$`)
func (uc *CustomerSessionUseCase) Update(
	ctx context.Context,
	id int32,
	req requests.UpdateCustomerSessionRequest,
) (responses.CustomerSessionUpdatedResponse, error) {

	if !reUsername.MatchString(req.TrainerUsername) {
		return responses.CustomerSessionUpdatedResponse{}, errors.New("invalid trainerUsername")
	}
	if req.PricePaid < 0 {
		return responses.CustomerSessionUpdatedResponse{}, errors.New("pricePaid must be >= 0")
	}
	if req.DiscountAmount < 0 {
		return responses.CustomerSessionUpdatedResponse{}, errors.New("discountAmount must be >= 0")
	}
	switch req.Status {
	case "ACTIVE", "EXPIRED", "COMPLETED", "CANCELLED":
	default:
		return responses.CustomerSessionUpdatedResponse{}, errors.New("invalid status")
	}

	// 2) trainer must exist & role=TRAINER
	n, err := uc.sessionRepo.CheckTrainerExists(ctx, req.TrainerUsername)
	if err != nil {
		return responses.CustomerSessionUpdatedResponse{}, err
	}
	if n == 0 {
		return responses.CustomerSessionUpdatedResponse{}, errors.New("trainer not found or not a TRAINER")
	}

	// 3) build params (decimal -> string)
	priceStr := fmt.Sprintf("%.2f", req.PricePaid)
	discStr := fmt.Sprintf("%.2f", req.DiscountAmount)

	params := repositories.UpdateCustomerSessionEditableFieldsParams{
		ID:              id,
		TrainerUsername: req.TrainerUsername,
		PricePaid:       priceStr,
		DiscountAmount:  sql.NullString{String: discStr, Valid: true},
		Status:          req.Status,
	}
	if err := uc.sessionRepo.UpdateEditableFields(ctx, params); err != nil {
		return responses.CustomerSessionUpdatedResponse{}, err
	}

	return responses.CustomerSessionUpdatedResponse{
		Message: fmt.Sprintf("Sessions Course ID: %d updated successfully", id),
	}, nil
}

func (uc *CustomerSessionUseCase) Delete(
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

func (uc *CustomerSessionUseCase) GetByID(ctx context.Context, id string) (responses.CustomerSession, error) {
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

	// --- แปลง DECIMAL(x,2) -> int64 ---
	pricePaid, err := utils.ParseDecimalToInt64(row.PricePaid, 2)
	if err != nil {
		return responses.CustomerSession{}, err
	}
	discount, err := utils.ParseDecimalToInt64(row.DiscountAmount.String, 2)
	if err != nil {
		return responses.CustomerSession{}, err
	}

	resp := responses.CustomerSession{
		ID:               utils.Itoa(row.ID),
		CustomerUsername: utils.NS(row.CustomerUsername),   // NullString -> string
		TrainerUsername:  utils.NS(row.TrainerUsername),    // NullString -> string
		SalesUsername:    utils.NS(row.SalesUsername),      // NullString -> string
		ProductID:        utils.Itoa(utils.NI32(row.ProductID)), // NullInt32 -> int32 -> string
		PurchaseDate:     utils.ToYMD(row.PurchaseDate),    // time.Time -> YYYY-MM-DD
		TotalSessions:    row.TotalSessions,                // int32
		UsedSessions:     utils.NI32(row.UsedSessions),     // NullInt32 -> int32
		PricePaid:        pricePaid,                        // int64
		DiscountAmount:   discount,                         // int64
		Status:           strings.ToUpper(string(row.Status)),
	}
	return resp, nil
}