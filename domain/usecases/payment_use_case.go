package usecases

import (
	"context"
	"fmt"
	"strconv"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/responses"
)

type PaymentUseCase struct {
	paymentRepo repositories.PaymentAccountRepository
}

func ProvidePaymentUseCase(paymentRepo repositories.PaymentAccountRepository) *PaymentUseCase {
	return &PaymentUseCase{
		paymentRepo: paymentRepo,
	}
}

// GetPaymentInfo - ดึงข้อมูลชำระเงินตาม Use Case 5S
// รับ productID และ discountAmount (optional) จาก frontend
func (u *PaymentUseCase) GetPaymentInfo(ctx context.Context, productID int32, discountAmount float64) (*responses.PaymentInfoResponse, error) {
	// Get payment info from database
	info, err := u.paymentRepo.GetPaymentInfoByProductId(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment info: %w", err)
	}

	// Parse list price from DECIMAL string
	listPrice, err := strconv.ParseFloat(info.ListPrice, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid list price: %w", err)
	}

	// Calculate payable amount
	payableAmount := listPrice - discountAmount
	if payableAmount < 0 {
		payableAmount = 0
	}

	// Build response
	response := &responses.PaymentInfoResponse{
		ProductID:        info.ProductID,
		ProductName:      info.ProductName,
		ProductType:      info.ProductType,
		ProductCategory:  info.ProductCategory,
		ListPrice:        listPrice,
		DiscountAmount:   discountAmount,
		PayableAmount:    payableAmount,
		PaymentAccountID: info.PaymentAccountID,
		AccountName:      info.AccountName,
		AccountNumber:    info.AccountNumber,
		BankName:         info.BankName,
		QRCodeURL:        info.QRCodeURL,
		AccountActive:    info.AccountActive,
	}

	// Add duration_days or session_amount based on product type
	if info.DurationDays.Valid {
		durationDays := info.DurationDays.Int32
		response.DurationDays = &durationDays
	}

	if info.SessionAmount.Valid {
		sessionAmount := info.SessionAmount.Int32
		response.SessionAmount = &sessionAmount
	}

	return response, nil
}
