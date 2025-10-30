package rest

import (
	"strconv"

	"github.com/WaritDev/private-fitness-backend/domain/usecases"
	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct {
	paymentUC *usecases.PaymentUseCase
}

func ProvidePaymentHandler(paymentUC *usecases.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{
		paymentUC: paymentUC,
	}
}

// GET /api/payments/info/:productId - ดึงข้อมูลชำระเงินตาม Use Case 5S
// Query params (optional): ?discountAmount=100.00
func (h *PaymentHandler) GetPaymentInfo(c *fiber.Ctx) error {
	// Get product ID from path
	productIDStr := c.Params("productId")
	productID, err := strconv.ParseInt(productIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid product ID",
			"result":      nil,
		})
	}

	// Get optional discount amount from query
	discountAmount := 0.0
	if discountStr := c.Query("discountAmount"); discountStr != "" {
		discount, err := strconv.ParseFloat(discountStr, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":      "Bad Request",
				"status_code": fiber.StatusBadRequest,
				"message":     "Invalid discount amount",
				"result":      nil,
			})
		}
		discountAmount = discount
	}

	// Call use case
	result, err := h.paymentUC.GetPaymentInfo(c.Context(), int32(productID), discountAmount)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":      "Not Found",
			"status_code": fiber.StatusNotFound,
			"message":     "Product not found or payment account not active",
			"result":      nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "OK",
		"status_code": fiber.StatusOK,
		"message":     "Payment info retrieved successfully",
		"result":      result,
	})
}
