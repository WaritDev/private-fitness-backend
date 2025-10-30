package usecases

import (
	"context"
	"fmt"
	"strconv"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/responses"
)

type ProductUseCase struct {
	repo repositories.ProductRepository
}

func ProvideProductUseCase(repo repositories.ProductRepository) *ProductUseCase {
	return &ProductUseCase{repo: repo}
}

// ListAllProducts returns all active products
func (u *ProductUseCase) ListAllProducts(ctx context.Context) ([]responses.ProductResponse, error) {
	products, err := u.repo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	result := make([]responses.ProductResponse, len(products))
	for i, p := range products {
		result[i] = u.mapToResponse(p)
	}
	return result, nil
}

// GetProductByID returns a single product by ID
func (u *ProductUseCase) GetProductByID(ctx context.Context, id int32) (responses.ProductResponse, error) {
	product, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return responses.ProductResponse{}, fmt.Errorf("product not found: %w", err)
	}

	return u.mapToResponse(product), nil
}

// ListDurations returns all DURATION type products
func (u *ProductUseCase) ListDurations(ctx context.Context) ([]responses.ProductResponse, error) {
	products, err := u.repo.ListDurations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list durations: %w", err)
	}

	result := make([]responses.ProductResponse, len(products))
	for i, p := range products {
		result[i] = u.mapToResponse(p)
	}
	return result, nil
}

// ListSessions returns all SESSION type products
func (u *ProductUseCase) ListSessions(ctx context.Context) ([]responses.ProductResponse, error) {
	products, err := u.repo.ListSessions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	result := make([]responses.ProductResponse, len(products))
	for i, p := range products {
		result[i] = u.mapToResponse(p)
	}
	return result, nil
}

// Helper to convert repository model to response
func (u *ProductUseCase) mapToResponse(p repositories.ProductInfo) responses.ProductResponse {
	// Parse list price from string to float64
	price, _ := strconv.ParseFloat(p.ListPrice, 64)

	return responses.ProductResponse{
		ID:               p.ID,
		Name:             p.Name,
		Type:             p.Type,
		Category:         p.Category,
		ListPrice:        price,
		DurationDays:     p.DurationDays,
		SessionAmount:    p.SessionAmount,
		IsActive:         p.IsActive == 1,
		PaymentAccountID: p.PaymentAccountID,
		CreatedAt:        p.CreatedAt,
		UpdatedAt:        p.UpdatedAt,
	}
}
