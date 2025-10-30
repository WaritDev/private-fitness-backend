package responses

import "database/sql"

// PaymentInfoResponse - Response สำหรับหน้าชำระเงิน (Use Case 5S)
type PaymentInfoResponse struct {
	// Product Information
	ProductID       int32   `json:"productId"`
	ProductName     string  `json:"productName"`
	ProductType     string  `json:"productType"`     // DURATION or SESSION
	ProductCategory string  `json:"productCategory"` // MUAY_THAI, KICK_BOXING, FITNESS, etc.
	ListPrice       float64 `json:"listPrice"`
	DurationDays    *int32  `json:"durationDays,omitempty"`  // สำหรับ DURATION
	SessionAmount   *int32  `json:"sessionAmount,omitempty"` // สำหรับ SESSION
	DiscountAmount  float64 `json:"discountAmount"`          // จาก frontend/calculation
	PayableAmount   float64 `json:"payableAmount"`           // ListPrice - DiscountAmount
	// Payment Account Information
	PaymentAccountID int32  `json:"paymentAccountId"`
	AccountName      string `json:"accountName"`
	AccountNumber    string `json:"accountNumber"`
	BankName         string `json:"bankName"`
	QRCodeURL        string `json:"qrCodeUrl"` // URL ของภาพ QR Code
	AccountActive    bool   `json:"accountActive"`
}

// PaymentInfoForCalculation - ข้อมูลจาก DB ก่อนคำนวณส่วนลด
type PaymentInfoForCalculation struct {
	ProductID        int32
	ProductName      string
	ProductType      string
	ProductCategory  string
	ListPrice        string // DECIMAL from DB
	DurationDays     sql.NullInt32
	SessionAmount    sql.NullInt32
	PaymentAccountID int32
	AccountName      string
	AccountNumber    string
	BankName         string
	QRCodeURL        string
	AccountActive    bool
}
