package repositories

import (
	"context"
	"time"

	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

type CustomerRepository interface {
	List(ctx context.Context) ([]dbmodel.ListCustomersRow, error)
	GetByUsername(ctx context.Context, username string) (dbmodel.GetCustomerByUsernameRow, error)
	CheckPhoneExistsExcept(ctx context.Context, phone, username string) (int64, error)
	CheckGmailExistsExcept(ctx context.Context, gmail, username string) (int64, error)

	UpdateUserNoPassword(ctx context.Context, p UpdateCustomerUserParams) error
	UpdateUserWithPassword(ctx context.Context, p UpdateCustomerUserWithPassParams) error
	UpdateCustomerDetail(ctx context.Context, p UpdateCustomerDetailParams) error

	DeleteCustomerCascade(ctx context.Context, username string) error
}

type UpdateCustomerUserParams struct {
	Username    string
	FirstName   string
	LastName    string
	Gender      string
	DateOfBirth time.Time
	PhoneNumber string
	Gmail       string
	IsActive    bool
}

type UpdateCustomerUserWithPassParams struct {
	UpdateCustomerUserParams
	HashedPass string
}

type UpdateCustomerDetailParams struct {
	Username                      string
	HealthInfo                    *string
	Address                       *string
	CompanyName                   *string
	CompanyPosition               *string
	MaritalStatus                 *string
	EmergencyContactName          *string
	EmergencyContactRelationship  *string
	EmergencyContactPhone         *string
	MarketingSource               *string
}