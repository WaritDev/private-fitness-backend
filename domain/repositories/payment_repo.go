package repositories

import (
	"context"

	"github.com/WaritDev/private-fitness-backend/domain/responses"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

// PaymentAccountRepository - Repository สำหรับจัดการข้อมูลการชำระเงิน
type PaymentAccountRepository interface {
	// GetPaymentInfoByProductId - ดึงข้อมูลสินค้าและบัญชีชำระเงิน (Q5S.1)
	GetPaymentInfoByProductId(ctx context.Context, productID int32) (*responses.PaymentInfoForCalculation, error)

	List(ctx context.Context, limit, offset int32) ([]dbmodel.ListPaymentAccountsRow, error)
	Count(ctx context.Context) (int64, error)

	Insert(ctx context.Context, p CreatePaymentAccountParams) (int32, error)

	Update(ctx context.Context, p UpdatePaymentAccountParams) error

	Delete(ctx context.Context, id int32) error

	// Payment Slip Verification Methods
	InsertPaymentVerification(ctx context.Context, params InsertPaymentVerificationParams) (int64, error)
	CheckDuplicatePayment(ctx context.Context, username string, productID int32, amount float64) (int64, error)
	UpdatePaymentVerificationStatus(ctx context.Context, params UpdatePaymentVerificationParams) error
	GetPaymentVerificationById(ctx context.Context, id int64) (*PaymentVerificationInfo, error)
}

type CreatePaymentAccountParams struct {
	AccountName   string
	AccountNumber string
	BankName      string
	QRCodeURL     string
	IsActive      *bool
}

type UpdatePaymentAccountParams struct {
	ID            int32
	AccountName   string
	AccountNumber string
	BankName      string
	QRCodeURL     string
	IsActive      *bool
}

// Payment Verification Structs
type InsertPaymentVerificationParams struct {
	CustomerUsername   string
	ProductID          int32
	Amount             float64
	SlipFilePath       string
	VerificationStatus string // "PENDING", "VERIFIED", "REJECTED"
}

type UpdatePaymentVerificationParams struct {
	ID                 int64
	VerificationStatus string
	SlipID             string
	Slip2GoResponse    string
}

type PaymentVerificationInfo struct {
	ID                 int64
	CustomerUsername   string
	ProductID          int32
	Amount             float64
	SlipFilePath       string
	SlipID             string
	VerificationStatus string
	Slip2GoResponse    string
	VerifiedAt         *string
	CreatedAt          string
}
