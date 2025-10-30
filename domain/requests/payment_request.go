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