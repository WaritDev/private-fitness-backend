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
