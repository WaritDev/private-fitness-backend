package repositories

import (
	"context"

	"github.com/WaritDev/private-fitness-backend/domain/responses"
)

// PaymentAccountRepository - Repository สำหรับจัดการข้อมูลการชำระเงิน
type PaymentAccountRepository interface {
	// GetPaymentInfoByProductId - ดึงข้อมูลสินค้าและบัญชีชำระเงิน (Q5S.1)
	GetPaymentInfoByProductId(ctx context.Context, productID int32) (*responses.PaymentInfoForCalculation, error)
}
