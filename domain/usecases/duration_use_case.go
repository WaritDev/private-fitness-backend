package usecases

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/responses"
)

type CustomerDurationUseCase struct {
	durationRepo repositories.CustomerDurationRepository
	productRepo  repositories.ProductRepository
}

func ProvideCustomerDurationUseCase(
	durationRepo repositories.CustomerDurationRepository,
	productRepo repositories.ProductRepository,
) *CustomerDurationUseCase {
	return &CustomerDurationUseCase{
		durationRepo: durationRepo,
		productRepo:  productRepo,
	}
}

// PurchaseDuration creates a new customer duration purchase
// Flow: ตรง use case description step 4
func (u *CustomerDurationUseCase) PurchaseDuration(ctx context.Context, username string, req requests.PurchaseDurationRequest) (int32, error) {
	// Get product details to calculate end date
	product, err := u.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return 0, fmt.Errorf("product not found: %w", err)
	}

	// Validate product type
	if product.Type != "DURATION" {
		return 0, fmt.Errorf("product is not a DURATION type")
	}

	// Validate duration days
	if product.DurationDays == nil || *product.DurationDays <= 0 {
		return 0, fmt.Errorf("invalid product: duration_days not set")
	}

	// Calculate dates
	now := time.Now()
	startDate := now
	endDate := startDate.AddDate(0, 0, int(*product.DurationDays))

	// Format prices as DECIMAL strings
	pricePaidStr := fmt.Sprintf("%.2f", req.PricePaid)
	discountAmountStr := fmt.Sprintf("%.2f", req.DiscountAmount)

	// Create duration record
	params := repositories.CreateCustomerDurationParams{
		CustomerUsername: username,
		SalesUsername:    req.SalesUsername,
		ProductID:        req.ProductID,
		PurchaseDate:     now,
		StartDate:        startDate,
		EndDate:          endDate,
		PricePaid:        pricePaidStr,
		DiscountAmount:   discountAmountStr,
		Status:           "ACTIVE",
	}

	if err := u.durationRepo.Create(ctx, params); err != nil {
		return 0, fmt.Errorf("failed to create duration: %w", err)
	}

	// Return success (in real implementation, we'd return the created ID)
	return product.ID, nil
}

// GetCustomerDurations returns all durations for a customer
func (u *CustomerDurationUseCase) GetCustomerDurations(ctx context.Context, username string) ([]responses.CustomerDurationResponse, error) {
	durations, err := u.durationRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer durations: %w", err)
	}

	result := make([]responses.CustomerDurationResponse, len(durations))
	for i, d := range durations {
		result[i] = u.mapToResponse(d)
	}
	return result, nil
}

// Helper to convert repository model to response
func (u *CustomerDurationUseCase) mapToResponse(d repositories.CustomerDurationInfo) responses.CustomerDurationResponse {
	pricePaid, _ := strconv.ParseFloat(d.PricePaid, 64)
	discountAmount, _ := strconv.ParseFloat(d.DiscountAmount, 64)

	return responses.CustomerDurationResponse{
		ID:               d.ID,
		CustomerUsername: d.CustomerUsername,
		SalesUsername:    d.SalesUsername,
		ProductID:        d.ProductID,
		PurchaseDate:     d.PurchaseDate,
		StartDate:        d.StartDate,
		EndDate:          d.EndDate,
		PricePaid:        pricePaid,
		DiscountAmount:   discountAmount,
		Status:           d.Status,
		CreatedAt:        d.CreatedAt,
		UpdatedAt:        d.UpdatedAt,
	}
}
