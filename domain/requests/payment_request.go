package requests

type ListPaymentAccountsRequest struct {
	Page  int32 `form:"page"  json:"page"`
	Limit int32 `form:"limit" json:"limit"`
}

type CreatePaymentAccountRequest struct {
	AccountName   string `json:"accountName"   binding:"required"`
	AccountNumber string `json:"accountNumber" binding:"required"`
	BankName      string `json:"bankName"      binding:"required"`
	QRCodeURL     string `json:"qrCodeUrl"     binding:"required"`
	IsActive      *bool  `json:"isActive"`
}

type UpdatePaymentAccountRequest struct {
	AccountName   string `json:"accountName"   binding:"required"`
	AccountNumber string `json:"accountNumber" binding:"required"`
	BankName      string `json:"bankName"      binding:"required"`
	QRCodeURL     string `json:"qrCodeUrl"     binding:"required,url"`
	IsActive      *bool  `json:"isActive"`
}

// VerifySlipPayload - Payload data sent from frontend (as JSON string in form-data)
type VerifySlipPayload struct {
	Username      string  `json:"username" validate:"required"`
	ProductID     int32   `json:"productId" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	AccountName   string  `json:"accountName" validate:"required"`
	AccountNumber string  `json:"accountNumber" validate:"required"`
	AccountType   string  `json:"accountType" validate:"required"` // e.g., "01004" for SCB
	PaymentDate   string  `json:"paymentDate,omitempty"`           // Optional: YYYY-MM-DD format
}
