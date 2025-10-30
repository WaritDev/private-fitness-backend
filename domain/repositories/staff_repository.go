package repositories

import (
	"context"
	"time"

	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

type StaffRepository interface {
	List(ctx context.Context, limit, offset int32) ([]dbmodel.ListStaffsRow, error)
	Count(ctx context.Context) (int64, error)
	CheckUsernameExists(ctx context.Context, username string) (int64, error)
	CheckPhoneExists(ctx context.Context, phone string) (int64, error)
	CheckGmailExists(ctx context.Context, gmail string) (int64, error)
	Create(ctx context.Context, p CreateStaffParams) error
	CheckPhoneExistsExceptUsername(ctx context.Context, phone, username string) (int64, error)
	CheckGmailExistsExceptUsername(ctx context.Context, gmail, username string) (int64, error)
	UpdateNoPassword(ctx context.Context, p UpdateStaffNoPassParams) error
	UpdateWithPassword(ctx context.Context, p UpdateStaffWithPassParams) error
	DeleteStaffCascade(ctx context.Context, username string) error
	GetByUsername(ctx context.Context, username string) (dbmodel.GetStaffByUsernameRow, error)
}

type CreateStaffParams struct {
	Username    string
	HashedPass  string
	Role        string
	FirstName   string
	LastName    string
	Gender      string
	DateOfBirth time.Time
	PhoneNumber string
	Gmail       string
	Specialty   *string
}

type UpdateStaffNoPassParams struct {
	Username    string
	Role        string
	FirstName   string
	LastName    string
	Gender      string
	DateOfBirth time.Time
	PhoneNumber string
	Gmail       string
	Specialty   *string
	IsActive    bool
}

type UpdateStaffWithPassParams struct {
	Username    string
	HashedPass  string
	Role        string
	FirstName   string
	LastName    string
	Gender      string
	DateOfBirth time.Time
	PhoneNumber string
	Gmail       string
	Specialty   *string
	IsActive    bool
}