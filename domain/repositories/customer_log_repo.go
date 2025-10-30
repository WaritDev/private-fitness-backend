package repositories

import (
	"context"
	"database/sql"
)

// CustomerLogRepository - Repository สำหรับจัดการ customer logs
type CustomerLogRepository interface {
	// CreateCustomerLog - สร้าง log สำหรับ customer activity
	CreateCustomerLog(ctx context.Context, tx *sql.Tx, customerUsername, logType string) error
}
