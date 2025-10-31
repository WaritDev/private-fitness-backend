package sql

import (
	"context"
	"database/sql"
	"strconv"

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

func (r *PaymentAccountRepository) List(ctx context.Context) ([]dbmodel.ListPaymentAccountsRow, error) {
	return r.q.ListPaymentAccounts(ctx)
}

func (r *PaymentAccountRepository) Insert(ctx context.Context, p repositories.CreatePaymentAccountParams) (int32, error) {
	isActive := utils.CoalesceBool(p.IsActive)
	res, err := r.q.InsertPaymentAccount(ctx, dbmodel.InsertPaymentAccountParams{
		AccountName:    p.AccountName,
		AccountNumber:  p.AccountNumber,
		BankName:       p.BankName,
		QrCodeImageUrl: p.QRCodeURL,
		IsActive:       isActive,
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
		AccountName:    p.AccountName,
		AccountNumber:  p.AccountNumber,
		BankName:       p.BankName,
		QrCodeImageUrl: p.QRCodeURL,
		IsActive:       isActive,
		ID:             p.ID,
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

// ========== Payment Slip Verification Methods ==========

// InsertPaymentVerification - Q_VERIFY_1: Insert payment verification log
func (r *PaymentAccountRepository) InsertPaymentVerification(ctx context.Context, params repositories.InsertPaymentVerificationParams) (int64, error) {
	res, err := r.q.InsertPaymentVerification(ctx, dbmodel.InsertPaymentVerificationParams{
		CustomerUsername:   params.CustomerUsername,
		ProductID:          params.ProductID,
		Amount:             utils.Decimal2(params.Amount),
		SlipFilePath:       sql.NullString{String: params.SlipFilePath, Valid: params.SlipFilePath != ""},
		VerificationStatus: dbmodel.PaymentVerificationsVerificationStatus(params.VerificationStatus),
	})
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// CheckDuplicatePayment - Q_VERIFY_2: Check duplicate payment slip
func (r *PaymentAccountRepository) CheckDuplicatePayment(ctx context.Context, username string, productID int32, amount float64) (int64, error) {
	count, err := r.q.CheckDuplicatePayment(ctx, dbmodel.CheckDuplicatePaymentParams{
		CustomerUsername: username,
		ProductID:        productID,
		Amount:           utils.Decimal2(amount),
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// UpdatePaymentVerificationStatus - Q_VERIFY_3: Update payment verification status after Slip2Go verification
func (r *PaymentAccountRepository) UpdatePaymentVerificationStatus(ctx context.Context, params repositories.UpdatePaymentVerificationParams) error {
	return r.q.UpdatePaymentVerificationStatus(ctx, dbmodel.UpdatePaymentVerificationStatusParams{
		VerificationStatus: dbmodel.PaymentVerificationsVerificationStatus(params.VerificationStatus),
		SlipID:             sql.NullString{String: params.SlipID, Valid: params.SlipID != ""},
		Slip2goResponse:    sql.NullString{String: params.Slip2GoResponse, Valid: params.Slip2GoResponse != ""},
		ID:                 int32(params.ID),
	})
}

// GetPaymentVerificationById - Q_VERIFY_4: Get payment verification by ID
func (r *PaymentAccountRepository) GetPaymentVerificationById(ctx context.Context, id int64) (*repositories.PaymentVerificationInfo, error) {
	row, err := r.q.GetPaymentVerificationById(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Parse amount from string to float64
	amount, _ := strconv.ParseFloat(row.Amount, 64)

	var verifiedAt *string
	if row.VerifiedAt.Valid {
		t := row.VerifiedAt.Time.Format("2006-01-02 15:04:05")
		verifiedAt = &t
	}

	var createdAt string
	if row.CreatedAt.Valid {
		createdAt = row.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}

	return &repositories.PaymentVerificationInfo{
		ID:                 int64(row.ID),
		CustomerUsername:   row.CustomerUsername,
		ProductID:          row.ProductID,
		Amount:             amount,
		SlipFilePath:       row.SlipFilePath.String,
		SlipID:             row.SlipID.String,
		VerificationStatus: string(row.VerificationStatus),
		Slip2GoResponse:    row.Slip2goResponse.String,
		VerifiedAt:         verifiedAt,
		CreatedAt:          createdAt,
	}, nil
}


func (r *PaymentAccountRepository) GetByID(ctx context.Context, id int32) (dbmodel.GetPaymentAccountByIDRow, error) {
	return r.q.GetPaymentAccountByID(ctx, id)
}