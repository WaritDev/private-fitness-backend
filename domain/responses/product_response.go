package responses

import "time"

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
