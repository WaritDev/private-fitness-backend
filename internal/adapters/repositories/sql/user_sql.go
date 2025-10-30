package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

type UserRepository struct {
	q  *dbmodel.Queries
	db *sql.DB
}

func ProvideUserRepository(q *dbmodel.Queries, db *sql.DB) *UserRepository {
	return &UserRepository{q: q, db: db}
}

func (r *UserRepository) List(ctx context.Context) ([]dbmodel.User, error) {
	rows, err := r.q.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	var users []dbmodel.User
	for _, row := range rows {
		users = append(users, dbmodel.User{
			Username:  row.Username,
			Gmail:     row.Gmail,
			CreatedAt: row.CreatedAt,
		})
	}
	return users, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (dbmodel.GetUserByUsernameRow, error) {
	return r.q.GetUserByUsername(ctx, username)
}

func (r *UserRepository) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	count, err := r.q.CheckUsernameExists(ctx, username)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CheckPhoneNumberExists checks if phone number exists in database (Q3S.1)
func (r *UserRepository) CheckPhoneNumberExists(ctx context.Context, phoneNumber string) (bool, error) {
	count, err := r.q.CheckPhoneNumberExists(ctx, phoneNumber)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CheckGmailExists checks if email exists in database (Q3S.2)
func (r *UserRepository) CheckGmailExists(ctx context.Context, email string) (bool, error) {
	count, err := r.q.CheckGmailExists(ctx, email)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) CreateUserWithCustomer(ctx context.Context, user repositories.CreateUserParams, customer repositories.CreateCustomerParams) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := r.q.WithTx(tx)

	// Parse date of birth
	var dobTime time.Time
	if user.DateOfBirth != nil && *user.DateOfBirth != "" {
		// Expecting format "YYYY-MM-DD"
		parsedTime, err := time.Parse("2006-01-02", *user.DateOfBirth)
		if err == nil {
			dobTime = parsedTime
		}
	}

	// Create user
	err = qtx.CreateUser(ctx, dbmodel.CreateUserParams{
		Username:    user.Username,
		Password:    user.Password,
		Gmail:       stringOrEmpty(user.Gmail),
		Role:        dbmodel.UsersRole(user.Role),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Gender:      dbmodel.UsersGender(stringOrEmpty(user.Gender)),
		DateOfBirth: dobTime,
		PhoneNumber: stringOrEmpty(user.PhoneNumber),
		Specialty:   sql.NullString{Valid: false},
		IsActive:    sql.NullBool{Bool: true, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Create customer
	err = qtx.CreateCustomer(ctx, dbmodel.CreateCustomerParams{
		Username:                     customer.Username,
		HealthInfo:                   stringOrEmpty(customer.HealthInfo),
		Address:                      stringOrEmpty(customer.Address),
		CompanyName:                  stringOrEmpty(customer.CompanyName),
		CompanyPosition:              stringOrEmpty(customer.CompanyPosition),
		MaritalStatus:                dbmodel.CustomersMaritalStatus(stringOrEmpty(customer.MaritalStatus)),
		EmergencyContactName:         stringOrEmpty(customer.EmergencyContactName),
		EmergencyContactRelationship: stringOrEmpty(customer.EmergencyContactRelationship),
		EmergencyContactPhone:        stringOrEmpty(customer.EmergencyContactPhone),
		MarketingSource:              stringOrEmpty(customer.MarketingSource),
	})
	if err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
