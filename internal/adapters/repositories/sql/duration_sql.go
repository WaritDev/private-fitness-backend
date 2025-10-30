package sql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

// CustomerDurationRepository implements domain/repositories.CustomerDurationRepository using sqlc
type CustomerDurationRepository struct {
	q *dbmodel.Queries
}

func ProvideCustomerDurationRepository(db *sql.DB) *CustomerDurationRepository {
	return &CustomerDurationRepository{q: dbmodel.New(db)}
}

func (r *CustomerDurationRepository) Create(ctx context.Context, params repositories.CreateCustomerDurationParams) error {
	var salesUsername sql.NullString
	if params.SalesUsername != nil {
		salesUsername = sql.NullString{String: *params.SalesUsername, Valid: true}
	}

	var discountAmount sql.NullString
	if params.DiscountAmount != "" {
		discountAmount = sql.NullString{String: params.DiscountAmount, Valid: true}
	}

	return r.q.CreateCustomerDuration(ctx, dbmodel.CreateCustomerDurationParams{
		CustomerUsername: sql.NullString{String: params.CustomerUsername, Valid: true},
		SalesUsername:    salesUsername,
		ProductID:        sql.NullInt32{Int32: params.ProductID, Valid: true},
		PurchaseDate:     params.PurchaseDate,
		StartDate:        params.StartDate,
		EndDate:          params.EndDate,
		PricePaid:        params.PricePaid,
		DiscountAmount:   discountAmount,
		Status:           dbmodel.CustomerDurationsStatus(params.Status),
	})
}

func (r *CustomerDurationRepository) GetByUsername(ctx context.Context, username string) ([]repositories.CustomerDurationInfo, error) {
	rows, err := r.q.GetCustomerDurationsByUsername(ctx, sql.NullString{String: username, Valid: true})
	if err != nil {
		return nil, err
	}

	result := make([]repositories.CustomerDurationInfo, len(rows))
	for i, row := range rows {
		result[i] = r.mapToCustomerDurationInfo(row)
	}
	return result, nil
}

func (r *CustomerDurationRepository) GetByID(ctx context.Context, id int32) (repositories.CustomerDurationInfo, error) {
	row, err := r.q.GetCustomerDurationById(ctx, id)
	if err != nil {
		return repositories.CustomerDurationInfo{}, err
	}

	return r.mapToCustomerDurationInfo(row), nil
}

// Helper function to map dbmodel.CustomerDuration to repositories.CustomerDurationInfo
func (r *CustomerDurationRepository) mapToCustomerDurationInfo(row dbmodel.CustomerDuration) repositories.CustomerDurationInfo {
	var salesUsername *string
	if row.SalesUsername.Valid {
		salesUsername = &row.SalesUsername.String
	}

	var customerUsername string
	if row.CustomerUsername.Valid {
		customerUsername = row.CustomerUsername.String
	}

	var productID int32
	if row.ProductID.Valid {
		productID = row.ProductID.Int32
	}

	var discountAmount string
	if row.DiscountAmount.Valid {
		discountAmount = row.DiscountAmount.String
	} else {
		discountAmount = "0.00"
	}

	var createdAt, updatedAt time.Time
	if row.CreatedAt.Valid {
		createdAt = row.CreatedAt.Time
	}
	if row.UpdatedAt.Valid {
		updatedAt = row.UpdatedAt.Time
	}

	return repositories.CustomerDurationInfo{
		ID:               row.ID,
		CustomerUsername: customerUsername,
		SalesUsername:    salesUsername,
		ProductID:        productID,
		PurchaseDate:     row.PurchaseDate,
		StartDate:        row.StartDate,
		EndDate:          row.EndDate,
		PricePaid:        row.PricePaid,
		DiscountAmount:   discountAmount,
		Status:           string(row.Status),
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}
}

func (r *CustomerDurationRepository) List(ctx context.Context, limit, offset int32) ([]dbmodel.ListCustomerDurationsRow, error) {
	return r.q.ListCustomerDurations(ctx, dbmodel.ListCustomerDurationsParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *CustomerDurationRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountCustomerDurations(ctx)
}

func (r *CustomerDurationRepository) GetDurationDaysForDurationID(ctx context.Context, durationID int32) (int32, error) {
	val, err := r.q.GetDurationDaysForDurationID(ctx, durationID)
	if err != nil {
		return 0, err
	}
	if !val.Valid {
		return 0, errors.New("duration_days is NULL")
	}
	return val.Int32, nil
}

func (r *CustomerDurationRepository) UpdateEditableFields(
	ctx context.Context,
	p repositories.UpdateCustomerDurationEditableFieldsParams,
) error {
	var disc sql.NullString
	if p.DiscountAmount != nil {
		disc = sql.NullString{String: *p.DiscountAmount, Valid: true}
	} else {
		disc = sql.NullString{Valid: false} // จะอัปเดตเป็น NULL
	}

	statusEnum := dbmodel.CustomerDurationsStatus(p.Status)
	return r.q.UpdateCustomerDurationEditableFields(ctx, dbmodel.UpdateCustomerDurationEditableFieldsParams{
		STRTODATE:      p.StartDateYMD,
		STRTODATE_2:    p.StartDateYMD,
		PricePaid:      p.PricePaid,
		DiscountAmount: disc,
		Status:         statusEnum,
		ID:             p.ID,
	})
}

func (r *CustomerDurationRepository) Delete(ctx context.Context, id int32) error {
	res, err := r.q.DeleteCustomerDurationByID(ctx, id)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}