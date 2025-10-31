package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

// CustomerLogRepository implements domain/repositories.CustomerLogRepository
type CustomerLogRepository struct {
	q *dbmodel.Queries
}

// ProvideCustomerLogRepository creates a new CustomerLogRepository
func ProvideCustomerLogRepository(q *dbmodel.Queries) repositories.CustomerLogRepository {
	return &CustomerLogRepository{q: q}
}

// CreateCustomerLog - สร้าง log สำหรับ customer activity
func (r *CustomerLogRepository) CreateCustomerLog(ctx context.Context, tx *sql.Tx, customerUsername, logType string) error {
	qtx := r.q.WithTx(tx)
	err := qtx.CreateCustomerLog(ctx, dbmodel.CreateCustomerLogParams{
		CustomerUsername: sql.NullString{String: customerUsername, Valid: true},
		LogType:          dbmodel.CustomerLogsLogType(logType),
	})
	if err != nil {
		return fmt.Errorf("failed to create customer log: %w", err)
	}
	return nil
}

func (r *CustomerLogRepository) List(ctx context.Context) ([]dbmodel.ListCustomerLogsRow, error) {
	return r.q.ListCustomerLogs(ctx)
}

func (r *CustomerLogRepository) UpdateByID(ctx context.Context, id int32, ts time.Time, logType string) (int64, error) {
	res, err := r.q.UpdateCustomerLogByID(ctx, dbmodel.UpdateCustomerLogByIDParams{
		CreatedAt: sql.NullTime{Time: ts, Valid: true},
		LogType:   dbmodel.CustomerLogsLogType(logType),
		ID:        id,
	})
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *CustomerLogRepository) DeleteByID(ctx context.Context, id int32) (int64, error) {
	res, err := r.q.DeleteCustomerLogByID(ctx, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *CustomerLogRepository) GetByID(ctx context.Context, id int32) (dbmodel.GetCustomerLogByIDRow, error) {
	return r.q.GetCustomerLogByID(ctx, id)
}

// CreatePendingCheckInLog - สร้าง pending check-in log
func (r *CustomerLogRepository) CreatePendingCheckInLog(ctx context.Context, customerUsername string, scheduleID int32) error {
	err := r.q.CreatePendingCheckInLog(ctx, dbmodel.CreatePendingCheckInLogParams{
		CustomerUsername: sql.NullString{String: customerUsername, Valid: true},
		ScheduleID:       sql.NullInt32{Int32: scheduleID, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to create pending check-in log: %w", err)
	}
	return nil
}

// GetPendingCheckInsByTrainer - ดึง pending check-ins ของ trainer
func (r *CustomerLogRepository) GetPendingCheckInsByTrainer(ctx context.Context, trainerUsername string) ([]repositories.PendingCheckInInfo, error) {
	rows, err := r.q.GetPendingCheckInsByTrainer(ctx, sql.NullString{String: trainerUsername, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get pending check-ins: %w", err)
	}

	result := make([]repositories.PendingCheckInInfo, len(rows))
	for i, row := range rows {
		var checkinTime time.Time
		if row.CheckinTime.Valid {
			checkinTime = row.CheckinTime.Time
		}

		result[i] = repositories.PendingCheckInInfo{
			LogID:                row.LogID,
			CustomerUsername:     row.CustomerUsername.String,
			CustomerFirstName:    row.CustomerFirstName,
			CustomerLastName:     row.CustomerLastName,
			CheckInTime:          checkinTime,
			ScheduleID:           row.ScheduleID,
			AppointmentStartTime: row.AppointmentStartTime,
			AppointmentEndTime:   row.AppointmentEndTime,
			SessionID:            row.SessionID,
			TotalSessions:        row.TotalSessions,
			UsedSessions:         row.UsedSessions,
		}
	}

	return result, nil
}

// UpdateCheckInLogStatus - อัปเดต status จาก PENDING เป็น CONFIRMED
func (r *CustomerLogRepository) UpdateCheckInLogStatus(ctx context.Context, logID int32) (int64, error) {
	res, err := r.q.UpdateCheckInLogStatus(ctx, logID)
	if err != nil {
		return 0, fmt.Errorf("failed to update check-in log status: %w", err)
	}
	return res.RowsAffected()
}
