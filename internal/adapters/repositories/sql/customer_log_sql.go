package sql

import (
	"context"
	"database/sql"
	"fmt"

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
