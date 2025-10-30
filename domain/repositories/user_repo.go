package repositories

import (
	"context"

	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

type CreateUserParams struct {
	Username    string
	Password    string
	Role        string
	FirstName   string
	LastName    string
	Gender      *string
	DateOfBirth *string
	PhoneNumber *string
	Gmail       *string
}

type CreateCustomerParams struct {
	Username                     string
	HealthInfo                   *string
	Address                      *string
	CompanyName                  *string
	CompanyPosition              *string
	MaritalStatus                *string
	EmergencyContactName         *string
	EmergencyContactRelationship *string
	EmergencyContactPhone        *string
	MarketingSource              *string
}

type UserRepo interface {
	List(ctx context.Context) ([]dbmodel.User, error)
	GetByUsername(ctx context.Context, username string) (dbmodel.GetUserByUsernameRow, error)
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
	CheckPhoneNumberExistsUser(ctx context.Context, phoneNumber string) (bool, error)
	CheckGmailExistsUser(ctx context.Context, email string) (bool, error)
	CreateUserWithCustomer(ctx context.Context, user CreateUserParams, customer CreateCustomerParams) error
	UpdateLoginTimestamp(ctx context.Context, username string) error // Q0S.2: Update updated_at on login
}
