package sql

import (
	"context"
	"database/sql"
	"time"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

// ProductRepository implements domain/repositories.ProductRepository using sqlc
type ProductRepository struct {
	q *dbmodel.Queries
}

// ProvideProductRepository creates a new ProductRepository
func ProvideProductRepository(q *dbmodel.Queries) repositories.ProductRepository {
	return &ProductRepository{q: q}
}

func (r *ProductRepository) ListAll(ctx context.Context) ([]repositories.ProductInfo, error) {
	rows, err := r.q.ListAllProducts(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]repositories.ProductInfo, len(rows))
	for i, row := range rows {
		// Use mapListRowToProductInfo for ListAllProductsRow type
		info := r.mapListRowToProductInfo(row.ID, row.Name, string(row.Type), string(row.Category), row.ListPrice, row.DurationDays, row.SessionAmount, row.IsActive, row.PaymentAccountID)
		// Add timestamps
		if row.CreatedAt.Valid {
			info.CreatedAt = row.CreatedAt.Time
		}
		if row.UpdatedAt.Valid {
			info.UpdatedAt = row.UpdatedAt.Time
		}
		result[i] = info
	}
	return result, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id int32) (repositories.ProductInfo, error) {
	row, err := r.q.GetProductById(ctx, id)
	if err != nil {
		return repositories.ProductInfo{}, err
	}

	return r.mapToProductInfo(row, true), nil
}

func (r *ProductRepository) ListDurations(ctx context.Context) ([]repositories.ProductInfo, error) {
	rows, err := r.q.ListDurations(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]repositories.ProductInfo, len(rows))
	for i, row := range rows {
		result[i] = r.mapListRowToProductInfo(row.ID, row.Name, string(row.Type), string(row.Category), row.ListPrice, row.DurationDays, row.SessionAmount, row.IsActive, row.PaymentAccountID)
	}
	return result, nil
}

func (r *ProductRepository) ListSessions(ctx context.Context) ([]repositories.ProductInfo, error) {
	rows, err := r.q.ListSessions(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]repositories.ProductInfo, len(rows))
	for i, row := range rows {
		result[i] = r.mapListRowToProductInfo(row.ID, row.Name, string(row.Type), string(row.Category), row.ListPrice, row.DurationDays, row.SessionAmount, row.IsActive, row.PaymentAccountID)
	}
	return result, nil
}

// Helper for List queries (without timestamps)
func (r *ProductRepository) mapListRowToProductInfo(id int32, name, typeStr, category, listPrice string, durationDays, sessionAmount sql.NullInt32, isActive sql.NullBool, paymentAccountID int32) repositories.ProductInfo {
	var durationDaysPtr, sessionAmountPtr *int32
	if durationDays.Valid {
		durationDaysPtr = &durationDays.Int32
	}
	if sessionAmount.Valid {
		sessionAmountPtr = &sessionAmount.Int32
	}

	var isActiveInt8 int8
	if isActive.Valid {
		if isActive.Bool {
			isActiveInt8 = 1
		}
	}

	return repositories.ProductInfo{
		ID:               id,
		Name:             name,
		Type:             typeStr,
		Category:         category,
		ListPrice:        listPrice,
		DurationDays:     durationDaysPtr,
		SessionAmount:    sessionAmountPtr,
		IsActive:         isActiveInt8,
		PaymentAccountID: paymentAccountID,
	}
}

// Helper function to convert dbmodel.Product to repositories.ProductInfo
func (r *ProductRepository) mapToProductInfo(row dbmodel.Product, includeTimestamps bool) repositories.ProductInfo {
	var durationDays, sessionAmount *int32
	if row.DurationDays.Valid {
		durationDays = &row.DurationDays.Int32
	}
	if row.SessionAmount.Valid {
		sessionAmount = &row.SessionAmount.Int32
	}

	var isActive int8
	if row.IsActive.Valid {
		if row.IsActive.Bool {
			isActive = 1
		} else {
			isActive = 0
		}
	}

	var createdAt, updatedAt time.Time
	if includeTimestamps {
		if row.CreatedAt.Valid {
			createdAt = row.CreatedAt.Time
		}
		if row.UpdatedAt.Valid {
			updatedAt = row.UpdatedAt.Time
		}
	}

	return repositories.ProductInfo{
		ID:               row.ID,
		Name:             row.Name,
		Type:             string(row.Type),
		Category:         string(row.Category),
		ListPrice:        row.ListPrice,
		DurationDays:     durationDays,
		SessionAmount:    sessionAmount,
		IsActive:         isActive,
		PaymentAccountID: row.PaymentAccountID,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}
}
