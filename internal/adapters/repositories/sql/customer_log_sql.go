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

func (r *CustomerLogRepository) List(ctx context.Context, limit, offset int32) ([]dbmodel.ListCustomerLogsRow, error) {
	return r.q.ListCustomerLogs(ctx, dbmodel.ListCustomerLogsParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *CustomerLogRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountCustomerLogs(ctx)
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