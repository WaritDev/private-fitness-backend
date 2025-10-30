package rest

import (
	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/usecases"
	"github.com/gofiber/fiber/v2"
)

type CustomerSessionHandler struct {
	useCase *usecases.CustomerSessionUseCase
}

func ProvideCustomerSessionHandler(useCase *usecases.CustomerSessionUseCase) *CustomerSessionHandler {
	return &CustomerSessionHandler{
		useCase: useCase,
	}
}

// Register - POST /api/customers/sessions/register
// Use Case 2.2C: ลงทะเบียนผู้ใช้งานสำหรับคอร์ส Sessions
func (h *CustomerSessionHandler) Register(c *fiber.Ctx) error {
	var req requests.RegisterCustomerSessionRequest

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
	result, err := h.useCase.RegisterCustomerSession(c.Context(), req)
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
		"message":     "Customer session registered successfully",
		"result":      result,
	})
}

// CheckPermission - GET /api/customers/sessions/check-permission
// ตรวจสอบสิทธิ์การเข้าถึงฟังก์ชันการจอง
func (h *CustomerSessionHandler) CheckPermission(c *fiber.Ctx) error {
	// Get customer username from query parameter
	customerUsername := c.Query("username")

	if customerUsername == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Username is required",
			"result":      nil,
		})
	}

	// Check permission
	hasPermission, err := h.useCase.CheckBookingPermission(c.Context(), customerUsername)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusInternalServerError,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	// Return result
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "success",
		"status_code": fiber.StatusOK,
		"message":     "Permission check completed",
		"result": fiber.Map{
			"hasPermission": hasPermission,
			"canBook":       hasPermission,
		},
	})
}

// GetActiveSessions - GET /api/customers/sessions/active/:username
// ดึงข้อมูล Session packages ที่ ACTIVE ของลูกค้า
func (h *CustomerSessionHandler) GetActiveSessions(c *fiber.Ctx) error {
	// Get customer username from path parameter
	customerUsername := c.Params("username")

	if customerUsername == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Username is required",
			"result":      nil,
		})
	}

	// Get active sessions
	sessions, err := h.useCase.GetCustomerActiveSessions(c.Context(), customerUsername)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusInternalServerError,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	// Success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "success",
		"status_code": fiber.StatusOK,
		"message":     "Active sessions retrieved successfully",
		"result":      sessions,
	})
}
