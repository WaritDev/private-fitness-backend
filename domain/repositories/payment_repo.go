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

	GetByID(ctx context.Context, id int32) (dbmodel.GetPaymentAccountByIDRow, error)
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