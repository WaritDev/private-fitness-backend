package sql

import (
	"context"
	"database/sql"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/responses"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

// PaymentAccountRepository implements domain/repositories.PaymentAccountRepository using sqlc
type PaymentAccountRepository struct {
	q *dbmodel.Queries
}

// ProvidePaymentAccountRepository creates a new PaymentAccountRepository
func ProvidePaymentAccountRepository(q *dbmodel.Queries) repositories.PaymentAccountRepository {
	return &PaymentAccountRepository{q: q}
}

// GetPaymentInfoByProductId - ดึงข้อมูลสินค้าและบัญชีชำระเงิน (Q5S.1)
func (r *PaymentAccountRepository) GetPaymentInfoByProductId(ctx context.Context, productID int32) (*responses.PaymentInfoForCalculation, error) {
	row, err := r.q.GetPaymentInfoByProductId(ctx, productID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err // Product not found or not active
		}
		return nil, err
	}

	return &responses.PaymentInfoForCalculation{
		ProductID:        row.ProductID,
		ProductName:      row.ProductName,
		ProductType:      string(row.ProductType),
		ProductCategory:  string(row.ProductCategory),
		ListPrice:        row.ListPrice,
		DurationDays:     row.DurationDays,
		SessionAmount:    row.SessionAmount,
		PaymentAccountID: row.PaymentAccountID,
		AccountName:      row.AccountName,
		AccountNumber:    row.AccountNumber,
		BankName:         row.BankName,
		QRCodeURL:        row.QrCodeImageUrl,
		AccountActive:    row.AccountActive,
	}, nil
}
