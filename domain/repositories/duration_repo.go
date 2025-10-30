package repositories

import (
	"context"
	"time"
)

type CustomerDurationRepository interface {
	Create(ctx context.Context, params CreateCustomerDurationParams) error
	GetByUsername(ctx context.Context, username string) ([]CustomerDurationInfo, error)
	GetByID(ctx context.Context, id int32) (CustomerDurationInfo, error)
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
