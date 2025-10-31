package responses

import (
	"time"

	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

type ListCustomerDurationsResponse struct {
	Data []dbmodel.ListCustomerDurationsRow `json:"data"`
	Meta PageMeta                           `json:"meta"`
}

// RegisterCustomerDurationResponse - Use Case 2.1C: Response after successful duration package registration
type RegisterCustomerDurationResponse struct {
	Username       string    `json:"username"`
	DurationID     int32     `json:"durationId"`
	ProductID      int32     `json:"productId"`
	SalesUsername  string    `json:"salesUsername"`
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
	DurationDays   int32     `json:"durationDays"`
	PricePaid      string    `json:"pricePaid"`
	DiscountAmount string    `json:"discountAmount"`
	Message        string    `json:"message"`
}

// CustomerDurationPackageResponse - Active duration package info (คล้าย CustomerSessionPackageResponse)
type CustomerDurationPackageResponse struct {
	ID               int32     `json:"id"`
	CustomerUsername string    `json:"customerUsername"`
	ProductID        int32     `json:"productId"`
	ProductName      string    `json:"productName"`  // JOIN from products
	DurationDays     int32     `json:"durationDays"` // จำนวนวันทั้งหมด
	SalesUsername    string    `json:"salesUsername"`
	PurchaseDate     time.Time `json:"purchaseDate"`
	StartDate        time.Time `json:"startDate"`
	EndDate          time.Time `json:"endDate"`
	DaysRemaining    int32     `json:"daysRemaining"`  // วันคงเหลือ (DATEDIFF)
	PricePaid        float64   `json:"pricePaid"`      // Convert from string to float64
	DiscountAmount   float64   `json:"discountAmount"` // Convert from string to float64
	Status           string    `json:"status"`         // ACTIVE, EXPIRED, CANCELLED
	CreatedAt        time.Time `json:"createdAt"`
}

// RenewDurationResponse - Response after customer self-purchase duration
type RenewDurationResponse struct {
	ID               int32     `json:"id"`
	CustomerUsername string    `json:"customerUsername"`
	ProductID        int32     `json:"productId"`
	ProductName      string    `json:"productName"`
	DurationDays     int32     `json:"durationDays"`
	SalesUsername    *string   `json:"salesUsername"` // NULL for self-purchase
	PurchaseDate     time.Time `json:"purchaseDate"`
	StartDate        time.Time `json:"startDate"`
	EndDate          time.Time `json:"endDate"`
	DaysRemaining    int32     `json:"daysRemaining"`
	PricePaid        float64   `json:"pricePaid"`
	DiscountAmount   float64   `json:"discountAmount"`
	Status           string    `json:"status"`
	Message          string    `json:"message"`
}
