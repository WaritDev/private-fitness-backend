package repositories

import (
	"context"
	"time"

	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

type ProductRepository interface {
	ListAll(ctx context.Context) ([]ProductInfo, error)
	GetByID(ctx context.Context, id int32) (ProductInfo, error)
	ListDurations(ctx context.Context) ([]ProductInfo, error)
	ListSessions(ctx context.Context) ([]ProductInfo, error)
	List(ctx context.Context) ([]dbmodel.ListProductsRow, error)
	CheckPaymentAccountActive(ctx context.Context, id int32) (int64, error)
	InsertDuration(ctx context.Context, p CreateProductDurationParams) (int32, error)
	InsertSession(ctx context.Context, p CreateProductSessionParams) (int32, error)
	UpdateDuration(ctx context.Context, p UpdateProductDurationParams) error
	UpdateSession(ctx context.Context, p UpdateProductSessionParams) error
	CountReferences(ctx context.Context, id int32) (int64, error)
    Delete(ctx context.Context, id int32) error
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

type CreateProductDurationParams struct {
	Name            string
	Category        string // ENUM (normalized)
	ListPrice       string // ส่งเป็น string ให้ลง DECIMAL
	DurationDays    int32
	IsActive        *bool
	PaymentAccountID int32
}

type CreateProductSessionParams struct {
	Name             string
	Category         string // ENUM (normalized)
	ListPrice        string // ส่งเป็น string ให้ลง DECIMAL
	SessionAmount    int32
	IsActive         *bool
	PaymentAccountID int32
}

type UpdateProductDurationParams struct {
	ID               int32
	Name             string
	Category         string
	ListPrice        string
	DurationDays     int32
	IsActive         *bool
	PaymentAccountID int32
}

type UpdateProductSessionParams struct {
	ID               int32
	Name             string
	Category         string
	ListPrice        string
	SessionAmount    int32
	IsActive         *bool
	PaymentAccountID int32
}

