package responses

import "time"

type CustomerDurationResponse struct {
	ID               int32     `json:"id"`
	CustomerUsername string    `json:"customerUsername"`
	SalesUsername    *string   `json:"salesUsername,omitempty"`
	ProductID        int32     `json:"productId"`
	PurchaseDate     time.Time `json:"purchaseDate"`
	StartDate        time.Time `json:"startDate"`
	EndDate          time.Time `json:"endDate"`
	PricePaid        float64   `json:"pricePaid"`
	DiscountAmount   float64   `json:"discountAmount"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type PurchaseSuccessResponse struct {
	DurationID int32  `json:"durationId"`
	Message    string `json:"message"`
}
type CustomerDurationUpdatedResponse struct {
	Message string `json:"message"`
}

type UpdateCustomerDurationRequest struct {
	StartDate      string   `json:"startDate"`
	PricePaid      float64  `json:"pricePaid"`
	DiscountAmount *float64 `json:"discountAmount"`
	Status         string   `json:"status"`
}

type CustomerDurationDeletedResponse struct {
	Message string `json:"message"`
}

type CustomerDuration struct {
	ID               string `json:"id"`
	CustomerUsername string `json:"customerUsername"`
	ProductID        string `json:"productId"`
	SalesUsername    string `json:"salesUsername"`
	PurchaseDate     string `json:"purchaseDate"` // YYYY-MM-DD
	StartDate        string `json:"startDate"`    // YYYY-MM-DD
	EndDate          string `json:"endDate"`      // YYYY-MM-DD
	PricePaid        int64  `json:"pricePaid"`
	DiscountAmount   int64  `json:"discountAmount"`
	Status           string `json:"status"`       // UPPERCASE (enum → string)
}