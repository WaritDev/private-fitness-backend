package requests

type PurchaseDurationRequest struct {
	ProductID      int32   `json:"productId" validate:"required"`
	SalesUsername  *string `json:"salesUsername"`
	PricePaid      float64 `json:"pricePaid" validate:"required,gt=0"`
	DiscountAmount float64 `json:"discountAmount" validate:"gte=0"`
}
