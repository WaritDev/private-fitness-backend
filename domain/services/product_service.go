package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/responses"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
	"github.com/WaritDev/private-fitness-backend/utils"
)

type ProductService struct {
	repo repositories.ProductRepository
}

func ProvideProductService(repo repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// ListAllProducts returns all active products
func (u *ProductService) ListAllProducts(ctx context.Context) ([]responses.ProductResponse, error) {
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
func (u *ProductService) GetProductByID(ctx context.Context, id int32) (responses.ProductResponse, error) {
	product, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return responses.ProductResponse{}, fmt.Errorf("product not found: %w", err)
	}

	return u.mapToResponse(product), nil
}

// ListDurations returns all DURATION type products
func (u *ProductService) ListDurations(ctx context.Context) ([]responses.ProductResponse, error) {
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
func (u *ProductService) ListSessions(ctx context.Context) ([]responses.ProductResponse, error) {
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
func (u *ProductService) mapToResponse(p repositories.ProductInfo) responses.ProductResponse {
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

func (uc *ProductService) List(ctx context.Context) ([]dbmodel.ListProductsRow, error) {
	return uc.repo.List(ctx)
}

// usecases/product_usecase.go
func (uc *ProductService) Create(ctx context.Context, req requests.CreateProductRequest) (responses.ProductCreatedResponse, error) {
    name := strings.TrimSpace(req.Name)
    if name == "" {
        return responses.ProductCreatedResponse{}, errors.New("name is required")
    }

    pt := normalizeType(req.ProductType) 
    cat := normalizeCategory(req.ProductCategory)

    if req.ListPrice < 0 {
        return responses.ProductCreatedResponse{}, errors.New("listPrice must be ≥ 0")
    }
    price2 := round2(req.ListPrice) 
    lp := fmt.Sprintf("%.2f", price2) 

    cnt, err := uc.repo.CheckPaymentAccountActive(ctx, req.PaymentAccountID)
    if err != nil {
        return responses.ProductCreatedResponse{}, err
    }
    if cnt == 0 {
        return responses.ProductCreatedResponse{}, errors.New("payment account not found or inactive")
    }

    switch pt {
    case "DURATION":
        if req.DurationDays == nil || *req.DurationDays <= 0 {
            return responses.ProductCreatedResponse{}, errors.New("durationDays must be > 0 for DURATION")
        }
        if req.SessionAmount != nil {
            return responses.ProductCreatedResponse{}, errors.New("sessionAmount must be null for DURATION")
        }

        id, err := uc.repo.InsertDuration(ctx, repositories.CreateProductDurationParams{
            Name:             name,
            Category:         cat,
            ListPrice:        lp,      
            DurationDays:     *req.DurationDays,
            IsActive:         req.IsActive,
            PaymentAccountID: req.PaymentAccountID,
        })
        if err != nil {
            return responses.ProductCreatedResponse{}, err
        }
        return responses.ProductCreatedResponse{
            ID:      id,
            Message: fmt.Sprintf("Product: %d created successfully", id),
        }, nil

    case "SESSION":
        if req.SessionAmount == nil || *req.SessionAmount <= 0 {
            return responses.ProductCreatedResponse{}, errors.New("sessionAmount must be > 0 for SESSION")
        }
        if req.DurationDays != nil {
            return responses.ProductCreatedResponse{}, errors.New("durationDays must be null for SESSION")
        }

        id, err := uc.repo.InsertSession(ctx, repositories.CreateProductSessionParams{
            Name:             name,
            Category:         cat,
            ListPrice:        lp,  
            SessionAmount:    *req.SessionAmount,
            IsActive:         req.IsActive,
            PaymentAccountID: req.PaymentAccountID,
        })
        if err != nil {
            return responses.ProductCreatedResponse{}, err
        }
        return responses.ProductCreatedResponse{
            ID:      id,
            Message: fmt.Sprintf("Product: %d created successfully", id),
        }, nil

    default:
        return responses.ProductCreatedResponse{}, errors.New("invalid product type")
    }
}

func normalizeType(s string) string {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "DURATION":
		return "DURATION"
	case "SESSION":
		return "SESSION"
	default:
		return ""
	}
}

func normalizeCategory(s string) string {
	x := strings.TrimSpace(strings.ToUpper(s))
	switch x {
	case "ECONOMY", "ECONOMIC":
		return "ECONOMIC"
	case "BUSINESS":
		return "BUSINESS"
	case "FIRST CLASS", "FIRST_CLASS", "FIRST-CLASS", "FIRSTCLASS", "FIRSTCLASSIC":
		return "FIRST_CLASS"
	default:
		return ""
	}
}

func (uc *ProductService) Update(
	ctx context.Context,
	id int32,
	req requests.UpdateProductRequest,
) (responses.ProductUpdatedResponse, error) {
	if req.ListPrice < 0 {
		return responses.ProductUpdatedResponse{}, errors.New("listPrice must be ≥ 0")
	}

	switch strings.ToUpper(req.Type) {
	case "DURATION":
		if req.DurationDays == nil || *req.DurationDays <= 0 {
			return responses.ProductUpdatedResponse{}, errors.New("durationDays must be > 0 for DURATION")
		}
		if err := uc.repo.UpdateDuration(ctx, repositories.UpdateProductDurationParams{
			ID:               id,
			Name:             req.Name,
			Category:         req.Category,
			ListPrice:        req.ListPrice,
			DurationDays:     *req.DurationDays,
			IsActive:         req.IsActive,
			PaymentAccountID: req.PaymentAccountID,
		}); err != nil {
			return responses.ProductUpdatedResponse{}, err
		}

	case "SESSION":
		if req.SessionAmount == nil || *req.SessionAmount <= 0 {
			return responses.ProductUpdatedResponse{}, errors.New("sessionAmount must be > 0 for SESSION")
		}
		if err := uc.repo.UpdateSession(ctx, repositories.UpdateProductSessionParams{
			ID:               id,
			Name:             req.Name,
			Category:         req.Category,
			ListPrice:        req.ListPrice,
			SessionAmount:    *req.SessionAmount,
			IsActive:         req.IsActive,
			PaymentAccountID: req.PaymentAccountID,
		}); err != nil {
			return responses.ProductUpdatedResponse{}, err
		}

	default:
		return responses.ProductUpdatedResponse{}, errors.New("invalid product type")
	}

	return responses.ProductUpdatedResponse{
		Message: fmt.Sprintf("Product: %d updated successfully", id),
	}, nil
}

func (uc *ProductService) Delete(
    ctx context.Context,
    id int32,
) (responses.ProductDeletedResponse, error) {

    refs, err := uc.repo.CountReferences(ctx, id)
    if err != nil {
        return responses.ProductDeletedResponse{}, err
    }

    if refs > 0 {
        return responses.ProductDeletedResponse{}, 
            errors.New("Cannot delete: Product is referenced by customer transactions.")
    }

    if err := uc.repo.Delete(ctx, id); err != nil {
        return responses.ProductDeletedResponse{}, err
    }

    return responses.ProductDeletedResponse{
        Message: fmt.Sprintf("Product: %d deleted successfully", id),
    }, nil
}

func (uc *ProductService) GetByID(ctx context.Context, id string) (responses.Product, error) {
	if strings.TrimSpace(id) == "" {
		return responses.Product{}, errors.New("id required")
	}

	id32, err := utils.Atoi32(id)
	if err != nil {
		return responses.Product{}, errors.New("invalid id")
	}

	row, err := uc.repo.GetByID(ctx, id32)
	if err != nil {
		return responses.Product{}, err
	}

	listPriceFloat, err := strconv.ParseFloat(row.ListPrice, 64)
	if err != nil {
		return responses.Product{}, fmt.Errorf("invalid list_price: %w", err)
	}
	listPriceFloat = round2(listPriceFloat)

	var duration int32
	if row.DurationDays != nil {
		duration = *row.DurationDays
	}

	var sessions int32
	if row.SessionAmount != nil {
		sessions = *row.SessionAmount
	}

	isActive := row.IsActive == 1

	createdAt := row.CreatedAt.Format(time.RFC3339)
	updatedAt := row.UpdatedAt.Format(time.RFC3339)

	resp := responses.Product{
		ID:               utils.Itoa(row.ID),
		Name:             row.Name,
		Type:             string(row.Type),
		Category:         string(row.Category),
		ListPrice:        listPriceFloat,
		DurationDays:     duration,
		SessionAmount:    sessions,
		IsActive:         isActive,
		PaymentAccountID: utils.Itoa(row.PaymentAccountID),
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}

	return resp, nil
}