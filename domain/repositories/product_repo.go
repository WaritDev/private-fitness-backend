package repositories

import (
	"context"
	"time"
)

type ProductRepository interface {
	ListAll(ctx context.Context) ([]ProductInfo, error)
	GetByID(ctx context.Context, id int32) (ProductInfo, error)
	ListDurations(ctx context.Context) ([]ProductInfo, error)
	ListSessions(ctx context.Context) ([]ProductInfo, error)
}

type ProductInfo struct {
	ID               int32
	Name             string
	Type             string
	Category         string
	ListPrice        string // DECIMAL from DB
	DurationDays     *int32
	SessionAmount    *int32
	IsActive         int8
	PaymentAccountID int32
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
