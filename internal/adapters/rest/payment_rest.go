package rest

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/WaritDev/private-fitness-backend/domain/requests"
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

func (h *PaymentHandler) List(c *fiber.Ctx) error {
	res, err := h.paymentUC.List(c.Context())
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusOK).JSON(res)
}

func (h *PaymentHandler) Create(c *fiber.Ctx) error {
	var req requests.CreatePaymentAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := h.paymentUC.Create(c.Context(), req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusOK).JSON(res)
}

// POST /api/payments/:id/update
func (h *PaymentHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var req requests.UpdatePaymentAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := h.paymentUC.Update(c.Context(), int32(id64), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// DELETE /api/payments/:id
func (h *PaymentHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	res, err := h.paymentUC.Delete(c.Context(), int32(id64))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

// POST /api/payment/verify-slip - Verify payment slip using Slip2Go API
func (h *PaymentHandler) VerifySlip(c *fiber.Ctx) error {
	// Step 1: Get file from multipart/form-data
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Missing slip image file",
		})
	}

	// Step 2: Open file to get io.Reader
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to open uploaded file",
		})
	}
	defer file.Close()

	// Step 3: Parse JSON payload from form data
	payloadStr := c.FormValue("payload")
	if payloadStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Missing payload data",
		})
	}

	var payload requests.VerifySlipPayload
	if err := c.App().Config().JSONDecoder([]byte(payloadStr), &payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid payload format",
		})
	}

	// Step 4: Call use case
	result, err := h.paymentUC.VerifySlip(c.Context(), payload, file, fileHeader.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Step 5: Return response based on verification result
	if result.Status == "error" {
		return c.Status(fiber.StatusBadRequest).JSON(result)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}


func (h *PaymentHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	resp, err := h.paymentUC.GetByID(c.Context(), id)
	if err != nil {
		if err.Error() == "payment account not found" || err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(resp)
}