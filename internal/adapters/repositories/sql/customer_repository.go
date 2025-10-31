package sql

import (
	"context"
	"database/sql"
	"strings"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
	"github.com/WaritDev/private-fitness-backend/utils"
)

type CustomerRepository struct {
	q  *dbmodel.Queries
	db *sql.DB
}

func ProvideCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{q: dbmodel.New(db), db: db}
}

func (r *CustomerRepository) List(ctx context.Context) ([]dbmodel.ListCustomersRow, error) {
	return r.q.ListCustomers(ctx)
}

func (r *CustomerRepository) GetByUsername(ctx context.Context, username string) (dbmodel.GetCustomerByUsernameRow, error) {
	return r.q.GetCustomerByUsername(ctx, username)
}

func (r *CustomerRepository) CheckPhoneExistsExcept(ctx context.Context, phone, username string) (int64, error) {
	return r.q.CheckCustomerPhoneExistsExcept(ctx, dbmodel.CheckCustomerPhoneExistsExceptParams{
		PhoneNumber: phone,
		Username:    username,
	})
}

func (r *CustomerRepository) CheckGmailExistsExcept(ctx context.Context, gmail, username string) (int64, error) {
	return r.q.CheckCustomerGmailExistsExcept(ctx, dbmodel.CheckCustomerGmailExistsExceptParams{
		LOWER:    strings.ToLower(gmail),
		Username: username,
	})
}

func (r *CustomerRepository) UpdateUserNoPassword(ctx context.Context, p repositories.UpdateCustomerUserParams) error {
	return r.q.UpdateCustomerUserNoPassword(ctx, dbmodel.UpdateCustomerUserNoPasswordParams{
		FirstName:   p.FirstName,
		LastName:    p.LastName,
		Gender:      dbmodel.UsersGender(p.Gender),
		DateOfBirth: p.DateOfBirth,
		PhoneNumber: p.PhoneNumber,
		Gmail:       strings.ToLower(p.Gmail),
		IsActive:    sql.NullBool{Bool: p.IsActive, Valid: true},
		Username:    p.Username,
	})
}

func (r *CustomerRepository) UpdateUserWithPassword(ctx context.Context, p repositories.UpdateCustomerUserWithPassParams) error {
	return r.q.UpdateCustomerUserWithPassword(ctx, dbmodel.UpdateCustomerUserWithPasswordParams{
		Password:    p.HashedPass,
		FirstName:   p.FirstName,
		LastName:    p.LastName,
		Gender:      dbmodel.UsersGender(p.Gender),
		DateOfBirth: p.DateOfBirth,
		PhoneNumber: p.PhoneNumber,
		Gmail:       strings.ToLower(p.Gmail),
		IsActive:    sql.NullBool{Bool: p.IsActive, Valid: true},
		Username:    p.Username,
	})
}

func (r *CustomerRepository) UpdateCustomerDetail(ctx context.Context, p repositories.UpdateCustomerDetailParams) error {
	return r.q.UpdateCustomersDetail(ctx, dbmodel.UpdateCustomersDetailParams{
		HealthInfo:                   utils.StrPtr(p.HealthInfo),
		Address:                      utils.StrPtr(p.Address),
		CompanyName:                  utils.StrPtr(p.CompanyName),
		CompanyPosition:              utils.StrPtr(p.CompanyPosition),
		MaritalStatus:                utils.EnumMarital(p.MaritalStatus),
		EmergencyContactName:         utils.StrPtr(p.EmergencyContactName),
		EmergencyContactRelationship: utils.StrPtr(p.EmergencyContactRelationship),
		EmergencyContactPhone:        utils.StrPtr(p.EmergencyContactPhone),
		MarketingSource:              utils.StrPtr(p.MarketingSource),
		Username:                     p.Username,
	})
}

func (r *CustomerRepository) DeleteCustomerCascade(ctx context.Context, username string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	qtx := r.q.WithTx(tx)

	ns := sql.NullString{String: username, Valid: true}

	// 1) CUSTOMER_LOG
	if err = qtx.DeleteCustomerLogsByCustomer(ctx, ns); err != nil {
		_ = tx.Rollback(); return err
	}
	// 2) TRAINING_SCHEDULE
	if err = qtx.DeleteTrainingScheduleByCustomer(ctx, ns); err != nil {
		_ = tx.Rollback(); return err
	}
	// 3) CUSTOMER_SESSION
	if err = qtx.DeleteCustomerSessionByCustomer(ctx, ns); err != nil {
		_ = tx.Rollback(); return err
	}
	// 4) CUSTOMER_DURATION
	if err = qtx.DeleteCustomerDurationByCustomer(ctx, ns); err != nil {
		_ = tx.Rollback(); return err
	}
	// 5) CUSTOMER
	if err = qtx.DeleteCustomerByUsername(ctx, username); err != nil {
		_ = tx.Rollback(); return err
	}
	// 6) USER
	if err = qtx.DeleteUserByUsername(ctx, username); err != nil {
		_ = tx.Rollback(); return err
	}

	return tx.Commit()
}