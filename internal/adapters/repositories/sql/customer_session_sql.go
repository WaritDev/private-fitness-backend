package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

// CustomerSessionRepository implements domain/repositories.CustomerSessionRepository
type CustomerSessionRepository struct {
	q  *dbmodel.Queries
	db *sql.DB
}

// ProvideCustomerSessionRepository creates a new CustomerSessionRepository
func ProvideCustomerSessionRepository(q *dbmodel.Queries, db *sql.DB) repositories.CustomerSessionRepository {
	return &CustomerSessionRepository{
		q:  q,
		db: db,
	}
}

// RegisterCustomerSession - Transaction: Create User → Customer → CustomerSession → TrainingSchedules → CustomerLog
func (r *CustomerSessionRepository) RegisterCustomerSession(ctx context.Context, tx *sql.Tx, params repositories.RegisterSessionParams) (int32, error) {
	// Use transaction queries
	qtx := r.q.WithTx(tx)

	// 1. Create User (Q2.2C.2 - part 1)
	err := qtx.CreateUser(ctx, dbmodel.CreateUserParams{
		Username:    params.Username,
		Password:    params.Password,
		Role:        dbmodel.UsersRoleCUSTOMER,
		FirstName:   params.FirstName,
		LastName:    params.LastName,
		Gender:      dbmodel.UsersGender(params.Gender),
		DateOfBirth: params.DateOfBirth,
		PhoneNumber: params.PhoneNumber,
		Gmail:       params.Gmail,
		Specialty:   sql.NullString{Valid: false},
		IsActive:    sql.NullBool{Bool: true, Valid: true},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	// 2. Create Customer (Q2.2C.2 - part 2)
	err = qtx.CreateCustomer(ctx, dbmodel.CreateCustomerParams{
		Username:                     params.Username,
		HealthInfo:                   params.HealthInfo,
		Address:                      params.Address,
		CompanyName:                  params.CompanyName,
		CompanyPosition:              params.CompanyPosition,
		MaritalStatus:                dbmodel.CustomersMaritalStatus(params.MaritalStatus),
		EmergencyContactName:         params.EmergencyContactName,
		EmergencyContactRelationship: params.EmergencyContactRelationship,
		EmergencyContactPhone:        params.EmergencyContactPhone,
		MarketingSource:              params.MarketingSource,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create customer: %w", err)
	}

	// 3. Create Customer Session (Q2.2C.2 - part 3)
	err = qtx.CreateCustomerSession(ctx, dbmodel.CreateCustomerSessionParams{
		CustomerUsername: sql.NullString{String: params.Username, Valid: true},
		TrainerUsername:  sql.NullString{String: params.TrainerUsername, Valid: true},
		ProductID:        sql.NullInt32{Int32: params.ProductID, Valid: true},
		SalesUsername:    sql.NullString{String: params.SalesUsername, Valid: true},
		PurchaseDate:     params.PurchaseDate,
		TotalSessions:    params.TotalSessions, // int32 - NOT NULL
		UsedSessions:     sql.NullInt32{Int32: params.UsedSessions, Valid: true},
		PricePaid:        params.PricePaid,
		DiscountAmount:   sql.NullString{String: params.DiscountAmount, Valid: true},
		Status:           dbmodel.CustomerSessionsStatusACTIVE,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create customer session: %w", err)
	}

	// Get the created session ID
	var sessionID int32
	err = tx.QueryRowContext(ctx, "SELECT LAST_INSERT_ID()").Scan(&sessionID)
	if err != nil {
		return 0, fmt.Errorf("failed to get session ID: %w", err)
	}

	// 4. Create Training Schedules (Q2C.3 - part 1)
	for _, schedule := range params.Schedules {
		err = qtx.CreateTrainingSchedule(ctx, dbmodel.CreateTrainingScheduleParams{
			TrainerUsername:  sql.NullString{String: schedule.TrainerUsername, Valid: true},
			CustomerUsername: sql.NullString{String: schedule.CustomerUsername, Valid: true},
			SessionID:        sql.NullInt32{Int32: sessionID, Valid: true},
			StartTime:        schedule.StartTime,
			EndTime:          schedule.EndTime,
			ScheduleType:     dbmodel.TrainingSchedulesScheduleTypeAPPOINTMENT,
		})
		if err != nil {
			return 0, fmt.Errorf("failed to create training schedule: %w", err)
		}
	}

	// 5. Create Customer Log (Q2C.3 - part 2)
	err = qtx.CreateCustomerLog(ctx, dbmodel.CreateCustomerLogParams{
		CustomerUsername: sql.NullString{String: params.Username, Valid: true},
		LogType:          dbmodel.CustomerLogsLogTypeBOOKSESSION,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create customer log: %w", err)
	}

	return sessionID, nil
}

// CheckBookingPermission - Q2C.1: ตรวจสอบสิทธิ์การเข้าถึงฟังก์ชันการจองก่อนโหลดปฏิทิน
// ตรวจสอบว่า Customer มีแพ็กเกจ Sessions แบบ ACTIVE หรือไม่
// หมายเหตุ: ถ้าทำครบแล้วจะเปลี่ยน status เป็น 'COMPLETED' โดยอัตโนมัติ
func (r *CustomerSessionRepository) CheckBookingPermission(ctx context.Context, customerUsername string) (bool, error) {
	count, err := r.q.CheckBookingPermission(ctx, sql.NullString{String: customerUsername, Valid: true})
	if err != nil {
		return false, fmt.Errorf("failed to check booking permission: %w", err)
	}

	// ถ้า count > 0 แสดงว่ามีแพ็กเกจ ACTIVE
	return count > 0, nil
}

// IncrementUsedSessions - Q3C.6: อัปเดตจำนวนครั้งที่ใช้ไป (used_sessions + 1)
func (r *CustomerSessionRepository) IncrementUsedSessions(ctx context.Context, tx *sql.Tx, sessionID int32) error {
	qtx := r.q.WithTx(tx)
	err := qtx.IncrementUsedSessions(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to increment used sessions: %w", err)
	}
	return nil
}

// GetActiveSessionByCustomer - หา Session package ACTIVE ของ Customer
func (r *CustomerSessionRepository) GetActiveSessionByCustomer(ctx context.Context, customerUsername string) (*repositories.ActiveSessionInfo, error) {
	row, err := r.q.GetActiveSessionByCustomer(ctx, sql.NullString{String: customerUsername, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No active session found
		}
		return nil, fmt.Errorf("failed to get active session: %w", err)
	}

	var trainerUsername string
	if row.TrainerUsername.Valid {
		trainerUsername = row.TrainerUsername.String
	}

	var customerUser string
	if row.CustomerUsername.Valid {
		customerUser = row.CustomerUsername.String
	}

	return &repositories.ActiveSessionInfo{
		ID:               row.ID,
		CustomerUsername: customerUser,
		TrainerUsername:  trainerUsername,
		TotalSessions:    row.TotalSessions,
		UsedSessions:     row.UsedSessions.Int32,
	}, nil
}

// DecrementUsedSessions - ยกเลิกการจอง: ลดจำนวนครั้งที่ใช้ไป (used_sessions - 1)
func (r *CustomerSessionRepository) DecrementUsedSessions(ctx context.Context, tx *sql.Tx, sessionID int32) error {
	qtx := r.q.WithTx(tx)
	err := qtx.DecrementUsedSessions(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to decrement used sessions: %w", err)
	}
	return nil
}

// IncrementUsedSessionsByUsername - Q5C.2: อัปเดตจำนวนครั้งที่ใช้ไปสำหรับ Check-in (Use Case 5C)
func (r *CustomerSessionRepository) IncrementUsedSessionsByUsername(ctx context.Context, customerUsername string) error {
	err := r.q.IncrementUsedSessionsByUsername(ctx, sql.NullString{String: customerUsername, Valid: true})
	if err != nil {
		return fmt.Errorf("failed to increment used sessions for %s: %w", customerUsername, err)
	}
	return nil
}

// GetCustomerActiveSessions - ดึงข้อมูล Session packages ที่ยัง ACTIVE ของลูกค้า
func (r *CustomerSessionRepository) GetCustomerActiveSessions(ctx context.Context, customerUsername string) ([]repositories.ActiveSessionPackageInfo, error) {
	rows, err := r.q.GetCustomerActiveSessions(ctx, sql.NullString{String: customerUsername, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get customer active sessions: %w", err)
	}

	result := make([]repositories.ActiveSessionPackageInfo, len(rows))
	for i, row := range rows {
		var trainerUsername string
		if row.TrainerUsername.Valid {
			trainerUsername = row.TrainerUsername.String
		}

		var customerUser string
		if row.CustomerUsername.Valid {
			customerUser = row.CustomerUsername.String
		}

		result[i] = repositories.ActiveSessionPackageInfo{
			ID:                row.ID,
			CustomerUsername:  customerUser,
			TrainerUsername:   trainerUsername,
			ProductID:         row.ProductID.Int32,
			ProductName:       row.ProductName,
			TotalSessions:     row.TotalSessions,
			UsedSessions:      row.UsedSessions.Int32,
			SessionsRemaining: row.SessionsRemaining,
			PurchaseDate:      row.PurchaseDate,
			PricePaid:         row.PricePaid,
			DiscountAmount:    row.DiscountAmount.String,
			Status:            string(row.Status),
		}
	}

	return result, nil
}

func (r *CustomerSessionRepository) List(ctx context.Context, limit, offset int32) ([]dbmodel.ListCustomerSessionsRow, error) {
	return r.q.ListCustomerSessions(ctx, dbmodel.ListCustomerSessionsParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *CustomerSessionRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountCustomerSessions(ctx)
}

func (r *CustomerSessionRepository) CheckTrainerExists(ctx context.Context, username string) (int64, error) {
	return r.q.CheckTrainerExists(ctx, username)
}

func (r *CustomerSessionRepository) UpdateEditableFields(
	ctx context.Context,
	p repositories.UpdateCustomerSessionEditableFieldsParams,
) error {
	return r.q.UpdateCustomerSessionEditableFields(ctx, dbmodel.UpdateCustomerSessionEditableFieldsParams{
		TrainerUsername: sql.NullString{String: p.TrainerUsername, Valid: true},
		PricePaid:       p.PricePaid,
		DiscountAmount:  p.DiscountAmount,
		Status:          dbmodel.CustomerSessionsStatus(p.Status),
		ID:              p.ID,
	})
}

func (r *CustomerSessionRepository) Delete(ctx context.Context, id int32) error {
	res, err := r.q.DeleteCustomerSessionByID(ctx, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *CustomerSessionRepository) GetByID(ctx context.Context, id int32) (dbmodel.GetCustomerSessionByIDRow, error) {
	return r.q.GetCustomerSessionByID(ctx, id)
}

// RenewSession - ต่ออายุ/ซื้อเพิ่ม Session Package (ลูกค้าซื้อเอง)
func (r *CustomerSessionRepository) RenewSession(ctx context.Context, params repositories.RenewSessionParams) error {
	return r.q.RenewCustomerSession(ctx, dbmodel.RenewCustomerSessionParams{
		CustomerUsername: sql.NullString{String: params.CustomerUsername, Valid: true},
		TrainerUsername:  sql.NullString{String: params.TrainerUsername, Valid: true},
		ProductID:        sql.NullInt32{Int32: params.ProductID, Valid: true},
		TotalSessions:    params.TotalSessions,
		PricePaid:        params.PricePaid,
	})
}
