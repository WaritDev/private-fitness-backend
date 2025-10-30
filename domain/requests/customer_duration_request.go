package requests

type ListCustomerDurationsRequest struct {
	Page  int32 `query:"page"`
	Limit int32 `query:"limit"`
}

type UpdateCustomerDurationRequest struct {
	StartDate      string  `json:"startDate"`       // YYYY-MM-DD
	PricePaid      float64 `json:"pricePaid"`       // >= 0
	DiscountAmount float64 `json:"discountAmount"`  // >= 0
	Status         string  `json:"status"`          // ACTIVE|EXPIRED|FROZEN|CANCELLED
}