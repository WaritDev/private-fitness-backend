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
}
