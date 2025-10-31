package usecases

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/WaritDev/private-fitness-backend/config"
	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/responses"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/slip2go"
	"github.com/WaritDev/private-fitness-backend/utils"
)

type PaymentUseCase struct {
	paymentRepo   repositories.PaymentAccountRepository
	slip2goClient *slip2go.Slip2GoClient
}

func ProvidePaymentUseCase(paymentRepo repositories.PaymentAccountRepository, cfg *config.Config) *PaymentUseCase {
	return &PaymentUseCase{
		paymentRepo:   paymentRepo,
		slip2goClient: slip2go.NewSlip2GoClient(cfg.Slip2GoSecretKey, cfg.MockSlip2Go),
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

func (uc *PaymentUseCase) List(ctx context.Context, req requests.ListPaymentAccountsRequest) (responses.ListPaymentAccountsResponse, error) {
	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	page := req.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	rows, err := uc.paymentRepo.List(ctx, limit, offset)
	if err != nil {
		return responses.ListPaymentAccountsResponse{}, err
	}
	total, err := uc.paymentRepo.Count(ctx)
	if err != nil {
		return responses.ListPaymentAccountsResponse{}, err
	}

	items := make([]responses.PaymentAccountItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, responses.PaymentAccountItem{
			ID:            r.ID,
			AccountName:   r.AccountName,
			AccountNumber: r.AccountNumber,
			BankName:      r.BankName,
			QRCodeURL:     r.QrCodeImageUrl,
			IsActive:      r.Column6,
		})
	}

	return responses.ListPaymentAccountsResponse{
		Data: items,
		Meta: responses.PageMeta{
			Page:       page,
			Limit:      limit,
			TotalItems: total,
			TotalPages: int32(math.Ceil(float64(total) / float64(limit))),
		},
	}, nil
}

var reURL = regexp.MustCompile(`^(https?://)[^\s]+$`)
var reAcctNo = regexp.MustCompile(`^[0-9\- ]{6,30}$`)

func (uc *PaymentUseCase) Create(ctx context.Context, req requests.CreatePaymentAccountRequest) (responses.PaymentAccountCreatedResponse, error) {
	if req.AccountName == "" {
		return responses.PaymentAccountCreatedResponse{}, errors.New("accountName cannot be empty")
	}
	if req.AccountNumber == "" || !reAcctNo.MatchString(req.AccountNumber) {
		return responses.PaymentAccountCreatedResponse{}, errors.New("invalid accountNumber")
	}
	if req.BankName == "" {
		return responses.PaymentAccountCreatedResponse{}, errors.New("bankName cannot be empty")
	}
	if req.QRCodeURL == "" || !reURL.MatchString(req.QRCodeURL) {
		return responses.PaymentAccountCreatedResponse{}, errors.New("invalid qrCodeUrl")
	}

	id, err := uc.paymentRepo.Insert(ctx, repositories.CreatePaymentAccountParams{
		AccountName:   req.AccountName,
		AccountNumber: req.AccountNumber,
		BankName:      req.BankName,
		QRCodeURL:     req.QRCodeURL,
		IsActive:      req.IsActive,
	})
	if err != nil {
		return responses.PaymentAccountCreatedResponse{}, err
	}

	return responses.PaymentAccountCreatedResponse{
		Message: "Payment Account created successfully",
		ID:      id,
	}, nil
}

func (uc *PaymentUseCase) Update(
	ctx context.Context,
	id int32,
	req requests.UpdatePaymentAccountRequest,
) (responses.PaymentAccountUpdatedResponse, error) {

	if req.AccountName == "" || req.AccountNumber == "" || req.BankName == "" {
		return responses.PaymentAccountUpdatedResponse{}, errors.New("missing required fields")
	}
	if req.QRCodeURL == "" || !reURL.MatchString(req.QRCodeURL) {
		return responses.PaymentAccountUpdatedResponse{}, errors.New("invalid qrCodeUrl")
	}

	err := uc.paymentRepo.Update(ctx, repositories.UpdatePaymentAccountParams{
		ID:            id,
		AccountName:   req.AccountName,
		AccountNumber: req.AccountNumber,
		BankName:      req.BankName,
		QRCodeURL:     req.QRCodeURL,
		IsActive:      req.IsActive,
	})
	if err != nil {
		return responses.PaymentAccountUpdatedResponse{}, err
	}

	return responses.PaymentAccountUpdatedResponse{
		Message: fmt.Sprintf("Payment Account: %d updated successfully", id),
	}, nil
}

func (uc *PaymentUseCase) Delete(ctx context.Context, id int32) (responses.PaymentAccountDeletedResponse, error) {
	if err := uc.paymentRepo.Delete(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			return responses.PaymentAccountDeletedResponse{}, fmt.Errorf("payment account not found")
		}
		return responses.PaymentAccountDeletedResponse{}, err
	}
	return responses.PaymentAccountDeletedResponse{
		Message: fmt.Sprintf("Payment Account: %d deleted successfully", id),
	}, nil
}

// ========== Payment Slip Verification ==========

// VerifySlip verifies payment slip using Slip2Go API (stateless - no database storage)
func (uc *PaymentUseCase) VerifySlip(ctx context.Context, payload requests.VerifySlipPayload, fileData io.Reader, filename string) (*responses.VerifySlipResponse, error) {
	// Call Slip2Go API to verify slip
	slip2goReq := slip2go.VerifySlipRequest{
		FileData:      fileData,
		Filename:      filename,
		AccountType:   payload.AccountType,
		AccountName:   payload.AccountName,
		AccountNumber: payload.AccountNumber,
		Amount:        payload.Amount,
		PaymentDate:   payload.PaymentDate,
	}

	slip2goResp, err := uc.slip2goClient.VerifySlip(slip2goReq)
	if err != nil {
		return &responses.VerifySlipResponse{
			Status:  "error",
			Message: fmt.Sprintf("Failed to verify slip: %s", err.Error()),
			Data: &responses.VerifySlipData{
				Verified: false,
			},
		}, nil
	}

	// Return verification result directly from Slip2Go
	if !slip2goResp.Result.Verified {
		return &responses.VerifySlipResponse{
			Status:  "error",
			Message: "Payment slip verification failed. Please check slip details and try again.",
			Data: &responses.VerifySlipData{
				SlipID:   slip2goResp.Result.SlipID,
				Verified: false,
			},
		}, nil
	}

	// Success - payment verified
	return &responses.VerifySlipResponse{
		Status:  "success",
		Message: "Payment verified successfully.",
		Data: &responses.VerifySlipData{
			SlipID:   slip2goResp.Result.SlipID,
			Verified: true,
		},
	}, nil
}

func (uc *PaymentUseCase) GetByID(ctx context.Context, paymentID string) (responses.PaymentAccount, error) {
	if strings.TrimSpace(paymentID) == "" {
		return responses.PaymentAccount{}, errors.New("paymentId required")
	}

	id32, err := utils.Atoi32(paymentID)
	if err != nil {
		return responses.PaymentAccount{}, errors.New("invalid paymentId")
	}

	row, err := uc.paymentRepo.GetByID(ctx, id32)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return responses.PaymentAccount{}, errors.New("payment account not found")
		}
		return responses.PaymentAccount{}, err
	}

	out := responses.PaymentAccount{
		ID:             utils.Itoa(row.ID),
		AccountName:    row.AccountName,
		AccountNumber:  row.AccountNumber,
		BankName:       row.BankName,
		QRCodeImageURL: row.QrCodeImageUrl,
		IsActive:       row.Column6,
	}
	return out, nil
}
