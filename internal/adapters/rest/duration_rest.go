package rest

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/usecases"
	"github.com/gofiber/fiber/v2"
)

type CustomerDurationHandler struct {
	UC     *usecases.CustomerDurationUseCase
	AuthUC *usecases.AuthUseCase
}

func ProvideCustomerDurationHandler(
	uc *usecases.CustomerDurationUseCase,
	authUC *usecases.AuthUseCase,
) *CustomerDurationHandler {
	return &CustomerDurationHandler{
		UC:     uc,
		AuthUC: authUC,
	}
}

// POST /api/durations/purchase
// ต้อง authenticated ก่อน (check JWT)
func (h *CustomerDurationHandler) PurchaseDuration(c *fiber.Ctx) error {
	// Extract and verify token
	token := c.Cookies("pf_auth")
	if token == "" {
		authHeader := c.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Unauthorized",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Authentication required",
			"result":      nil,
		})
	}

	// Verify token
	payload, err := h.AuthUC.VerifyToken(c.Context(), token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Unauthorized",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Invalid or expired token",
			"result":      nil,
		})
	}

	// Parse request
	var req requests.PurchaseDurationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid request body",
			"result":      nil,
		})
	}

	// Validate required fields
	if req.ProductID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "Product ID is required",
			"result":      nil,
		})
	}

	if req.PricePaid <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "Price paid must be greater than 0",
			"result":      nil,
		})
	}

	// Purchase duration
	_, err = h.UC.PurchaseDuration(c.Context(), payload.Sub, req)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = fiber.StatusNotFound
		} else if strings.Contains(err.Error(), "invalid") {
			statusCode = fiber.StatusBadRequest
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"status":      "Error",
			"status_code": statusCode,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "Created",
		"status_code": fiber.StatusCreated,
		"message":     "Package purchased successfully",
		"result":      fiber.Map{"ok": true},
	})
}

// GET /api/durations/my
// Get all durations for current authenticated user
func (h *CustomerDurationHandler) GetMyDurations(c *fiber.Ctx) error {
	// Extract and verify token
	token := c.Cookies("pf_auth")
	if token == "" {
		authHeader := c.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Unauthorized",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Authentication required",
			"result":      nil,
		})
	}

	// Verify token
	payload, err := h.AuthUC.VerifyToken(c.Context(), token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "Unauthorized",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Invalid or expired token",
			"result":      nil,
		})
	}

	// Get durations
	durations, err := h.UC.GetCustomerDurations(c.Context(), payload.Sub)
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
		"message":     "Durations retrieved successfully",
		"result":      durations,
	})
}

// GetActiveDuration - GET /api/customers/durations/active/:username
// ดึงข้อมูล Duration packages ที่ ACTIVE ของลูกค้า (คล้าย GetActiveSessions)
func (h *CustomerDurationHandler) GetActiveDuration(c *fiber.Ctx) error {
	username := c.Params("username")
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "Username is required",
			"result":      nil,
		})
	}

	// เรียก use case
	packages, err := h.UC.GetCustomerActiveDuration(c.Context(), username)
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
		"message":     "Active duration packages retrieved successfully",
		"result":      packages,
	})
}

func (h *CustomerDurationHandler) ListDurations(c *fiber.Ctx) error {
	resp, err := h.UC.List(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *CustomerDurationHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	var body requests.UpdateCustomerDurationRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	resp, err := h.UC.UpdateDuration(c.Context(), int32(id64), body)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *CustomerDurationHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	resp, err := h.UC.Delete(c.Context(), int32(id64))
	if err != nil {
		if err.Error() == "duration package not found" {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// Register - POST /api/customers/durations/register
// Use Case 2.1C: ลงทะเบียนผู้ใช้งานสำหรับแพ็กเกจ Duration
func (h *CustomerDurationHandler) Register(c *fiber.Ctx) error {
	var req requests.RegisterCustomerDurationRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid request body",
			"result":      nil,
		})
	}

	// Validate required fields
	if req.Username == "" || req.Password == "" || req.ConfirmPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Username, password, and confirm password are required",
			"result":      nil,
		})
	}

	// Check password match
	if req.Password != req.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Passwords do not match",
			"result":      nil,
		})
	}

	// Call use case
	result, err := h.UC.RegisterCustomerDuration(c.Context(), req)
	if err != nil {
		// Check for specific error types
		errMsg := err.Error()
		if errMsg == "USERNAME_ALREADY_EXISTS" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":      "error",
				"status_code": fiber.StatusBadRequest,
				"message":     "Username already exists",
				"result":      nil,
			})
		}

		// Generic error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusInternalServerError,
			"message":     errMsg,
			"result":      nil,
		})
	}

	// Success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "success",
		"status_code": fiber.StatusOK,
		"message":     "Customer duration registered successfully",
		"result":      result,
	})
}

func (h *CustomerDurationHandler) GetByID(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "id required"})
	}

	out, err := h.UC.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "customer_duration not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(out)
}