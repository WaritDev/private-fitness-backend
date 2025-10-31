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
	q  *dbmodel.Queries
	db *sql.DB
}

// ProvideCustomerDurationRepository creates a new CustomerDurationRepository
func ProvideCustomerDurationRepository(q *dbmodel.Queries, db *sql.DB) repositories.CustomerDurationRepository {
	return &CustomerDurationRepository{
		q:  q,
		db: db,
	}
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

// RegisterDuration - Use Case 2.1C: Transaction: Create User → Customer → CustomerDuration
func (r *CustomerDurationRepository) RegisterDuration(ctx context.Context, tx *sql.Tx, params repositories.RegisterDurationParams) (int32, error) {
	// Use transaction queries
	qtx := r.q.WithTx(tx)

	// 1. Create User (Q2.1C.2 - part 1)
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
		return 0, err
	}

	// 2. Create Customer (Q2.1C.2 - part 2)
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
		return 0, err
	}

	// 3. Create Customer Duration (Q2.1C.2 - part 3)
	err = qtx.CreateCustomerDuration(ctx, dbmodel.CreateCustomerDurationParams{
		CustomerUsername: sql.NullString{String: params.Username, Valid: true},
		SalesUsername:    sql.NullString{String: params.SalesUsername, Valid: true},
		ProductID:        sql.NullInt32{Int32: params.ProductID, Valid: true},
		PurchaseDate:     params.PurchaseDate,
		StartDate:        params.StartDate,
		EndDate:          params.EndDate,
		PricePaid:        params.PricePaid,
		DiscountAmount:   sql.NullString{String: params.DiscountAmount, Valid: true},
		Status:           dbmodel.CustomerDurationsStatusACTIVE,
	})
	if err != nil {
		return 0, err
	}

	// Get the created duration ID
	var durationID int32
	err = tx.QueryRowContext(ctx, "SELECT LAST_INSERT_ID()").Scan(&durationID)
	if err != nil {
		return 0, err
	}

	return durationID, nil
}

// GetCustomerActiveDuration - ดึง Duration packages ที่ยัง ACTIVE ของลูกค้า (คล้าย GetCustomerActiveSessions)
func (r *CustomerDurationRepository) GetCustomerActiveDuration(ctx context.Context, username string) ([]repositories.ActiveDurationPackageInfo, error) {
	rows, err := r.q.GetCustomerActiveDuration(ctx, sql.NullString{String: username, Valid: true})
	if err != nil {
		return nil, err
	}

	result := make([]repositories.ActiveDurationPackageInfo, 0, len(rows))
	for _, row := range rows {
		// Handle NULL values from JOIN
		customerUsername := ""
		if row.CustomerUsername.Valid {
			customerUsername = row.CustomerUsername.String
		}

		productID := int32(0)
		if row.ProductID.Valid {
			productID = row.ProductID.Int32
		}

		durationDays := int32(0)
		if row.DurationDays.Valid {
			durationDays = row.DurationDays.Int32
		}

		salesUsername := ""
		if row.SalesUsername.Valid {
			salesUsername = row.SalesUsername.String
		}

		discountAmount := ""
		if row.DiscountAmount.Valid {
			discountAmount = row.DiscountAmount.String
		}

		createdAt := time.Time{}
		if row.CreatedAt.Valid {
			createdAt = row.CreatedAt.Time
		}

		result = append(result, repositories.ActiveDurationPackageInfo{
			ID:               row.ID,
			CustomerUsername: customerUsername,
			ProductID:        productID,
			ProductName:      row.ProductName, // string (not NULL from JOIN)
			DurationDays:     durationDays,
			SalesUsername:    salesUsername,
			PurchaseDate:     row.PurchaseDate,
			StartDate:        row.StartDate,
			EndDate:          row.EndDate,
			DaysRemaining:    row.DaysRemaining, // int32 (calculated, not NULL)
			PricePaid:        row.PricePaid,     // string (DECIMAL as string)
			DiscountAmount:   discountAmount,
			Status:           string(row.Status),
			CreatedAt:        createdAt,
		})
	}

	return result, nil
}

// RenewDuration - Customer self-purchase duration package (ซื้อเพิ่มเอง)
func (r *CustomerDurationRepository) RenewDuration(ctx context.Context, params repositories.RenewDurationParams) error {
	return r.q.RenewCustomerDuration(ctx, dbmodel.RenewCustomerDurationParams{
		CustomerUsername: sql.NullString{String: params.CustomerUsername, Valid: true},
		ProductID:        sql.NullInt32{Int32: params.ProductID, Valid: true},
		PricePaid:        params.PricePaid,
		ID:               params.ProductID, // WHERE p.id = ? (validate product)
	})
}
