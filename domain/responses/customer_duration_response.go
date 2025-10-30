package responses

import "github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
import "time"

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
