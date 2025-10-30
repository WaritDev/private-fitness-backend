package repositories

import (
	"context"
	"database/sql"
	"time"
)

type CustomerDurationRepository interface {
	Create(ctx context.Context, params CreateCustomerDurationParams) error
	GetByUsername(ctx context.Context, username string) ([]CustomerDurationInfo, error)
	GetByID(ctx context.Context, id int32) (CustomerDurationInfo, error)
	// RegisterDuration - Use Case 2.1C: สร้าง User, Customer, CustomerDuration ในครั้งเดียว (Transaction)
	RegisterDuration(ctx context.Context, tx *sql.Tx, params RegisterDurationParams) (int32, error)
}

type CreateCustomerDurationParams struct {
	CustomerUsername string
	SalesUsername    *string
	ProductID        int32
	PurchaseDate     time.Time
	StartDate        time.Time
	EndDate          time.Time
	PricePaid        string // DECIMAL as string
	DiscountAmount   string // DECIMAL as string
	Status           string
}

type CustomerDurationInfo struct {
	ID               int32
	CustomerUsername string
	SalesUsername    *string
	ProductID        int32
	PurchaseDate     time.Time
	StartDate        time.Time
	EndDate          time.Time
	PricePaid        string
	DiscountAmount   string
	Status           string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// RegisterDurationParams - Parameters สำหรับ RegisterDuration (Use Case 2.1C)
type RegisterDurationParams struct {
	// User table
	Username    string
	Password    string // Already hashed
	FirstName   string
	LastName    string
	Gender      string
	DateOfBirth time.Time
	PhoneNumber string
	Gmail       string

	// Customer table
	HealthInfo                   string
	Address                      string
	CompanyName                  string
	CompanyPosition              string
	MaritalStatus                string
	EmergencyContactName         string
	EmergencyContactRelationship string
	EmergencyContactPhone        string
	MarketingSource              string

	// Customer Duration table
	ProductID      int32
	SalesUsername  string
	PurchaseDate   time.Time
	StartDate      time.Time
	EndDate        time.Time
	PricePaid      string // DECIMAL as string
	DiscountAmount string // DECIMAL as string
}
