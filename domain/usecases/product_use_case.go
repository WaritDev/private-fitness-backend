package usecases

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/requests"
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

func (uc *ProductUseCase) List(ctx context.Context, req requests.ListProductsRequest) (responses.ListProductsResponse, error) {
	limit := req.Limit
	if limit <= 0 || limit > 100 { limit = 10 }
	page := req.Page
	if page <= 0 { page = 1 }
	offset := (page - 1) * limit

	data, err := uc.repo.List(ctx, limit, offset)
	if err != nil { return responses.ListProductsResponse{}, err }

	total, err := uc.repo.Count(ctx)
	if err != nil { return responses.ListProductsResponse{}, err }

	return responses.ListProductsResponse{
		Data: data,
		Meta: responses.PageMeta{
			Page:       page,
			Limit:      limit,
			TotalItems: total,
			TotalPages: int32(math.Ceil(float64(total)/float64(limit))),
		},
	}, nil
}

func (uc *ProductUseCase) Create(ctx context.Context, req requests.CreateProductRequest) (responses.ProductCreatedResponse, error) {
	// 1) Validate base
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return responses.ProductCreatedResponse{}, errors.New("name is required")
	}

	pt := normalizeType(req.ProductType)        // DURATION|SESSION
	cat := normalizeCategory(req.ProductCategory) // ECONOMY|BUSINESS|FIRST_CLASS

	if pt == "" {
		return responses.ProductCreatedResponse{}, errors.New("invalid productType")
	}
	if cat == "" {
		return responses.ProductCreatedResponse{}, errors.New("invalid productCategory")
	}
	if req.ListPrice < 0 {
		return responses.ProductCreatedResponse{}, errors.New("listPrice must be >= 0")
	}

	// 2) Check Payment Account active
	cnt, err := uc.repo.CheckPaymentAccountActive(ctx, req.PaymentAccountID)
	if err != nil {
		return responses.ProductCreatedResponse{}, err
	}
	if cnt == 0 {
		return responses.ProductCreatedResponse{}, errors.New("payment account not found or inactive")
	}

	// 3) Branch by type + validate cross fields
	lp := fmt.Sprintf("%.2f", req.ListPrice)

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
	}

	return responses.ProductCreatedResponse{}, errors.New("unreachable")
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

func (uc *ProductUseCase) Update(
    ctx context.Context,
    id int32,
    req requests.UpdateProductRequest,
) (responses.ProductUpdatedResponse, error) {

    if req.Name == "" {
        return responses.ProductUpdatedResponse{}, errors.New("name cannot be empty")
    }
    if req.ListPrice < 0 {
        return responses.ProductUpdatedResponse{}, errors.New("listPrice must be >= 0")
    }
    if req.Category == "" {
        return responses.ProductUpdatedResponse{}, errors.New("category is required")
    }

    // ตรวจสอบเงื่อนไขตามประเภท
    switch req.Type {
    case "DURATION":
        if req.DurationDays == nil || *req.DurationDays <= 0 {
            return responses.ProductUpdatedResponse{}, errors.New("durationDays must be > 0 for DURATION")
        }
        priceStr := fmt.Sprintf("%.2f", req.ListPrice)
		err := uc.repo.UpdateDuration(ctx, repositories.UpdateProductDurationParams{
			ID:               id,
			Name:             req.Name,
			Category:         req.Category,
			ListPrice:        priceStr,
			DurationDays:     *req.DurationDays,
			IsActive:         req.IsActive,
			PaymentAccountID: req.PaymentAccountID,
		})
        if err != nil {
            return responses.ProductUpdatedResponse{}, err
        }

    case "SESSION":
        if req.SessionAmount == nil || *req.SessionAmount <= 0 {
            return responses.ProductUpdatedResponse{}, errors.New("sessionAmount must be > 0 for SESSION")
        }
		priceStr := fmt.Sprintf("%.2f", req.ListPrice)
        err := uc.repo.UpdateSession(ctx, repositories.UpdateProductSessionParams{
            ID:              id,
            Name:            req.Name,
            Category:        req.Category,
            ListPrice:       priceStr,
            SessionAmount:   *req.SessionAmount,
            IsActive:        req.IsActive,
            PaymentAccountID: req.PaymentAccountID,
        })
        if err != nil {
            return responses.ProductUpdatedResponse{}, err
        }

    default:
        return responses.ProductUpdatedResponse{}, errors.New("invalid product type")
    }

    return responses.ProductUpdatedResponse{
        Message: fmt.Sprintf("Product: %d updated successfully", id),
    }, nil
}

func (uc *ProductUseCase) Delete(
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