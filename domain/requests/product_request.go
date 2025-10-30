package requests

type ProductFilterRequest struct {
	Type     string `json:"type"`     // DURATION or SESSION
	Category string `json:"category"` // ECONOMIC, BUSINESS, FIRST_CLASS
}

type ListProductsRequest struct {
	Page  int32 `query:"page"`
	Limit int32 `query:"limit"`
}

type CreateProductRequest struct {
	Name             string   `json:"name"`
	ProductType      string   `json:"productType"`      // DURATION | SESSION
	ProductCategory  string   `json:"productCategory"`  // ECONOMY|BUSINESS|FIRST_CLASS (รองรับสะกดแบบ UI แล้ว normalize)
	ListPrice        float64  `json:"listPrice"`
	DurationDays     *int32   `json:"durationDays,omitempty"`
	SessionAmount    *int32   `json:"sessionAmount,omitempty"`
	IsActive         *bool    `json:"isActive,omitempty"`
	PaymentAccountID int32    `json:"paymentAccountId"`
}

type UpdateProductRequest struct {
    Name             string   `json:"name"`
    Type             string   `json:"type"`
    Category         string   `json:"category"`
    ListPrice        float64  `json:"listPrice"`
    DurationDays     *int32   `json:"durationDays,omitempty"`
    SessionAmount    *int32   `json:"sessionAmount,omitempty"`
    IsActive         *bool    `json:"isActive,omitempty"`
    PaymentAccountID int32    `json:"paymentAccountId"`
}