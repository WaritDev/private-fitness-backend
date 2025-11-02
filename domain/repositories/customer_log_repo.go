package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

// CustomerLogRepository - Repository สำหรับจัดการ customer logs
type CustomerLogRepository interface {
	// CreateCustomerLog - สร้าง log สำหรับ customer activity
	CreateCustomerLog(ctx context.Context, tx *sql.Tx, customerUsername, logType string) error

	List(ctx context.Context) ([]dbmodel.ListCustomerLogsRow, error)

	UpdateByID(ctx context.Context, id int32, ts time.Time, logType string) (int64, error)
	DeleteByID(ctx context.Context, id int32) (int64, error)

	GetByID(ctx context.Context, id int32) (dbmodel.GetCustomerLogByIDRow, error)

	// Check-in Flow: Pending and Confirmed
	// CreatePendingCheckInLog - สร้าง pending check-in log
	CreatePendingCheckInLog(ctx context.Context, customerUsername string, scheduleID int32) error

	// GetPendingCheckInsByTrainer - ดึง pending check-ins ของ trainer
	GetPendingCheckInsByTrainer(ctx context.Context, trainerUsername string) ([]PendingCheckInInfo, error)

	// UpdateCheckInLogStatus - อัปเดต status จาก PENDING เป็น CONFIRMED
	UpdateCheckInLogStatus(ctx context.Context, tx *sql.Tx, logID int32) (int64, error)
}

// PendingCheckInInfo - ข้อมูล pending check-in
type PendingCheckInInfo struct {
	LogID              int32
	CustomerUsername   string
	CustomerFirstName  string
	CustomerLastName   string
	CheckInTime        time.Time
	ScheduleID         int32
	AppointmentStartTime time.Time
	AppointmentEndTime   time.Time
	SessionID          sql.NullInt32
	TotalSessions      sql.NullInt32
	UsedSessions       sql.NullInt32
}
