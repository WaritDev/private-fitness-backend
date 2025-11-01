package responses

import (
	"time"

	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

type ProductResponse struct {
	ID               int32     `json:"id"`
	Name             string    `json:"name"`
	Type             string    `json:"type"`
	Category         string    `json:"category"`
	ListPrice        float64   `json:"listPrice"`
	DurationDays     *int32    `json:"durationDays,omitempty"`
	SessionAmount    *int32    `json:"sessionAmount,omitempty"`
	IsActive         bool      `json:"isActive"`
	PaymentAccountID int32     `json:"paymentAccountId"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type ListProductsResponse struct {
	Data []dbmodel.ListProductsRow `json:"data"`
	Meta PageMeta                  `json:"meta"`
}

type ProductCreatedResponse struct {
	ID      int32  `json:"id"`
	Message string `json:"message"`
}

type ProductUpdatedResponse struct {
	Message string `json:"message"`
}

type ProductDeletedResponse struct {
    Message string `json:"message"`
}

type Product struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	Category         string `json:"category"`
	ListPrice        float64  `json:"listPrice"`
	DurationDays     int32  `json:"durationDays,omitempty"`
	SessionAmount    int32  `json:"sessionAmount,omitempty"`
	IsActive         bool   `json:"isActive"`
	PaymentAccountID string `json:"paymentAccountId"`
	CreatedAt        string `json:"createdAt,omitempty"`
	UpdatedAt        string `json:"updatedAt,omitempty"`
}