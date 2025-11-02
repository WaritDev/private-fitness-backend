package sql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
	"github.com/WaritDev/private-fitness-backend/utils"
)

type StaffRepository struct {
	q  *dbmodel.Queries
	db *sql.DB
}

func ProvideStaffRepository(db *sql.DB) *StaffRepository {
	return &StaffRepository{
		q:  dbmodel.New(db),
		db: db,
	}
}

func (r *StaffRepository) List(ctx context.Context) ([]dbmodel.ListStaffsRow, error) {
	return r.q.ListStaffs(ctx)
}

func (r *StaffRepository) CheckUsernameExists(ctx context.Context, username string) (int64, error) {
	return r.q.CheckUsernameExists(ctx, username)
}
func (r *StaffRepository) CheckPhoneExists(ctx context.Context, phone string) (int64, error) {
	return r.q.CheckPhoneExists(ctx, phone)
}
func (r *StaffRepository) CheckGmailExists(ctx context.Context, gmail string) (int64, error) {
	return r.q.CheckGmailExists(ctx, gmail)
}
func (r *StaffRepository) Create(ctx context.Context, p repositories.CreateStaffParams) error {
	return r.q.CreateStaff(ctx, dbmodel.CreateStaffParams{
		Username:    p.Username,
		Password:    p.HashedPass,
		Role:        dbmodel.UsersRole(p.Role),
		FirstName:   p.FirstName,
		LastName:    p.LastName,
		Gender:      dbmodel.UsersGender(p.Gender),
		DateOfBirth: p.DateOfBirth,
		PhoneNumber: p.PhoneNumber,
		Gmail:       p.Gmail,
		Specialty:   utils.ToNullString(p.Specialty),
	})
}

func (r *StaffRepository) GetByUsername(ctx context.Context, username string) (dbmodel.GetStaffByUsernameRow, error) {
	return r.q.GetStaffByUsername(ctx, username)
}
func (r *StaffRepository) CheckPhoneExistsExceptUsername(ctx context.Context, phone, username string) (int64, error) {
	return r.q.CheckPhoneExistsExceptUsername(ctx, dbmodel.CheckPhoneExistsExceptUsernameParams{
		PhoneNumber: phone,
		Username:    username,
	})
}

func (r *StaffRepository) CheckGmailExistsExceptUsername(ctx context.Context, gmail, username string) (int64, error) {
	return r.q.CheckGmailExistsExceptUsername(ctx, dbmodel.CheckGmailExistsExceptUsernameParams{
		Gmail:    gmail,
		Username: username,
	})
}

func (r *StaffRepository) UpdateNoPassword(ctx context.Context, p repositories.UpdateStaffNoPassParams) error {
	var ns sql.NullString
	if p.Specialty != nil {
		ns = sql.NullString{String: *p.Specialty, Valid: true}
	}

	return r.q.UpdateStaffNoPassword(ctx, dbmodel.UpdateStaffNoPasswordParams{
		Role:        dbmodel.UsersRole(p.Role),
		FirstName:   p.FirstName,
		LastName:    p.LastName,
		Gender:      dbmodel.UsersGender(p.Gender),
		DateOfBirth: p.DateOfBirth,
		PhoneNumber: p.PhoneNumber,
		Gmail:       p.Gmail,
		Specialty:   ns,
		IsActive:    sql.NullBool{Bool: p.IsActive, Valid: true},
		Username:    p.Username,
	})
}

func (r *StaffRepository) UpdateWithPassword(ctx context.Context, p repositories.UpdateStaffWithPassParams) error {
	var ns sql.NullString
	if p.Specialty != nil {
		ns = sql.NullString{String: *p.Specialty, Valid: true}
	}

	return r.q.UpdateStaffWithPassword(ctx, dbmodel.UpdateStaffWithPasswordParams{
		Password:    p.HashedPass,
		Role:        dbmodel.UsersRole(p.Role),
		FirstName:   p.FirstName,
		LastName:    p.LastName,
		Gender:      dbmodel.UsersGender(p.Gender),
		DateOfBirth: p.DateOfBirth,
		PhoneNumber: p.PhoneNumber,
		Gmail:       p.Gmail,
		Specialty:   ns,
		IsActive:    sql.NullBool{Bool: p.IsActive, Valid: true},
		Username:    p.Username,
	})
}

func (r *StaffRepository) DeleteStaffCascade(ctx context.Context, username string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	qtx := r.q.WithTx(tx)

	role, err := qtx.GetUserRole(ctx, username)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	ns := sql.NullString{String: username, Valid: true}

	switch role {
	case dbmodel.UsersRoleTRAINER:
		if err = qtx.DeleteTrainerAvailabilityByTrainer(ctx, username); err != nil { _ = tx.Rollback(); return err }
		if err = qtx.DeleteTrainingScheduleByTrainer(ctx, ns); err != nil { _ = tx.Rollback(); return err }
		if err = qtx.DeleteCustomerSessionByTrainer(ctx, ns); err != nil { _ = tx.Rollback(); return err }

	case dbmodel.UsersRoleSALES:
		if err = qtx.DeleteCustomerDurationBySales(ctx, ns); err != nil { _ = tx.Rollback(); return err }
		if err = qtx.DeleteCustomerSessionBySales(ctx, ns); err != nil { _ = tx.Rollback(); return err }

	case dbmodel.UsersRoleADMIN, dbmodel.UsersRoleMANAGER, dbmodel.UsersRoleCUSTOMER:

	default:
		_ = tx.Rollback()
		return errors.New("unknown role")
	}

	if err = qtx.DeleteUserByUsername(ctx, username); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}