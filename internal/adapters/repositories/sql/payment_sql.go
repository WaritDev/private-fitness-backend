package sql

import (
	"context"
	"database/sql"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/responses"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
	"github.com/WaritDev/private-fitness-backend/utils"
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

func (r *PaymentAccountRepository) List(ctx context.Context, limit, offset int32) ([]dbmodel.ListPaymentAccountsRow, error) {
	return r.q.ListPaymentAccounts(ctx, dbmodel.ListPaymentAccountsParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *PaymentAccountRepository) Count(ctx context.Context) (int64, error) {
	return r.q.CountPaymentAccounts(ctx)
}

func (r *PaymentAccountRepository) Insert(ctx context.Context, p repositories.CreatePaymentAccountParams) (int32, error) {
	isActive := utils.CoalesceBool(p.IsActive)
	res, err := r.q.InsertPaymentAccount(ctx, dbmodel.InsertPaymentAccountParams{
		AccountName:     p.AccountName,
		AccountNumber:   p.AccountNumber,
		BankName:        p.BankName,
		QrCodeImageUrl:  p.QRCodeURL,
		IsActive:        isActive,
	})
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int32(id), nil
}

func (r *PaymentAccountRepository) Update(ctx context.Context, p repositories.UpdatePaymentAccountParams) error {
	isActive := utils.CoalesceBool(p.IsActive)
	return r.q.UpdatePaymentAccountByID(ctx, dbmodel.UpdatePaymentAccountByIDParams{
		AccountName:     p.AccountName,
		AccountNumber:   p.AccountNumber,
		BankName:        p.BankName,
		QrCodeImageUrl:  p.QRCodeURL,
		IsActive:        isActive,
		ID:              p.ID,
	})
}

func (r *PaymentAccountRepository) Delete(ctx context.Context, id int32) error {
    res, err := r.q.DeletePaymentAccountByID(ctx, id)
    if err != nil {
        return err
    }
    n, err := res.RowsAffected()
    if err == nil && n == 0 {
        return sql.ErrNoRows
    }
    return err
}