package repositories

import (
	"context"
	"time"

	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

type CustomerDurationRepository interface {
	Create(ctx context.Context, params CreateCustomerDurationParams) error
	GetByUsername(ctx context.Context, username string) ([]CustomerDurationInfo, error)
	GetByID(ctx context.Context, id int32) (CustomerDurationInfo, error)

	List(ctx context.Context, limit, offset int32) ([]dbmodel.ListCustomerDurationsRow, error)
	Count(ctx context.Context) (int64, error)

	GetDurationDaysForDurationID(ctx context.Context, durationID int32) (int32, error)
	UpdateEditableFields(ctx context.Context, p UpdateCustomerDurationEditableFieldsParams) error

	Delete(ctx context.Context, id int32) error
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

type CustomerDurationTableRow struct {
	DurationID          int32
	CustomerUsername    string
	CustomerFirstName   string
	CustomerLastName    string
	ProductID           int32
	ProductName         string
	ProductType         string
	ProductCategory     string
	DurationDays        int32
	SalesUsername       string
	PurchaseDate        time.Time
	StartDate           time.Time
	EndDate             time.Time
	PricePaid           string
	DiscountAmount      string
	Status              string
}

type UpdateCustomerDurationEditableFieldsParams struct {
	ID             int32   // cd.id
	StartDateYMD   string  // "YYYY-MM-DD"
	PricePaid      string  // e.g. "5000.00"
	DiscountAmount *string // nil = ไม่อัปเดต/ให้เป็น NULL, ไม่ nil = "0.00", "100.00"
	Status         string  // "ACTIVE" | "EXPIRED" | "FROZEN" | "CANCELLED"
}