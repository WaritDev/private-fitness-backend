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
