package rest

import (
	"strconv"

	"github.com/WaritDev/private-fitness-backend/domain/usecases"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	UC *usecases.ProductUseCase
}

func ProvideProductHandler(uc *usecases.ProductUseCase) *ProductHandler {
	return &ProductHandler{UC: uc}
}

// GET /api/products
func (h *ProductHandler) ListAllProducts(c *fiber.Ctx) error {
	products, err := h.UC.ListAllProducts(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "Internal Server Error",
			"status_code": fiber.StatusInternalServerError,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "OK",
		"status_code": fiber.StatusOK,
		"message":     "Products retrieved successfully",
		"result":      products,
	})
}

// GET /api/products/:id
func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid product ID",
			"result":      nil,
		})
	}

	product, err := h.UC.GetProductByID(c.Context(), int32(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":      "Not Found",
			"status_code": fiber.StatusNotFound,
			"message":     "Product not found",
			"result":      nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "OK",
		"status_code": fiber.StatusOK,
		"message":     "Product retrieved successfully",
		"result":      product,
	})
}

// GET /api/durations
func (h *ProductHandler) ListDurations(c *fiber.Ctx) error {
	products, err := h.UC.ListDurations(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "Internal Server Error",
			"status_code": fiber.StatusInternalServerError,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "OK",
		"status_code": fiber.StatusOK,
		"message":     "Duration packages retrieved successfully",
		"result":      products,
	})
}

// GET /api/sessions
func (h *ProductHandler) ListSessions(c *fiber.Ctx) error {
	products, err := h.UC.ListSessions(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "Internal Server Error",
			"status_code": fiber.StatusInternalServerError,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "OK",
		"status_code": fiber.StatusOK,
		"message":     "Session packages retrieved successfully",
		"result":      products,
	})
}
