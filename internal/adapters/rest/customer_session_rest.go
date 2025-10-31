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

type CustomerSessionHandler struct {
	useCase *usecases.CustomerSessionUseCase
	authUC  *usecases.AuthUseCase
}

func ProvideCustomerSessionHandler(
	useCase *usecases.CustomerSessionUseCase,
	authUC *usecases.AuthUseCase,
) *CustomerSessionHandler {
	return &CustomerSessionHandler{
		useCase: useCase,
		authUC:  authUC,
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
// Q2C.1: ตรวจสอบสิทธิ์การเข้าถึงฟังก์ชันการจองก่อนโหลดปฏิทิน
// ตรวจสอบว่า Customer มีแพ็กเกจ Sessions แบบ ACTIVE หรือไม่
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

func (h *CustomerSessionHandler) ListSessions(c *fiber.Ctx) error {
	resp, err := h.useCase.List(c.Context(), requests.ListCustomerSessionsRequest{})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *CustomerSessionHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	var req requests.UpdateCustomerSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	resp, err := h.useCase.Update(c.Context(), int32(id64), req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *CustomerSessionHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	res, err := h.useCase.Delete(c.Context(), int32(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *CustomerSessionHandler) GetByID(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Params("id"))
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "id required"})
	}

	resp, err := h.useCase.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "customer_session not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.JSON(resp)
}

// RenewSession - POST /api/customer-sessions/renew
// Use Case: ต่ออายุ/ซื้อเพิ่ม Session Package (ลูกค้าซื้อเอง)
func (h *CustomerSessionHandler) RenewSession(c *fiber.Ctx) error {
	// Step 1: Extract JWT token from Cookie or Authorization header
	token := c.Cookies("pf_auth")
	if token == "" {
		authHeader := c.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Unauthorized: Please login first",
			"result":      nil,
		})
	}

	// Step 2: Verify token and extract username
	payload, err := h.authUC.VerifyToken(c.Context(), token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusUnauthorized,
			"message":     "Invalid or expired token",
			"result":      nil,
		})
	}

	customerUsername := payload.Sub // username from JWT payload

	var req requests.RenewSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid request body",
			"result":      nil,
		})
	}

	// Validate required fields
	if req.ProductID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Product ID is required and must be greater than 0",
			"result":      nil,
		})
	}

	if req.TrainerUsername == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Trainer username is required",
			"result":      nil,
		})
	}

	// Step 3: Call use case
	result, err := h.useCase.RenewSession(c.Context(), customerUsername, req)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		message := err.Error()

		// Handle specific errors
		if strings.Contains(message, "PRODUCT_NOT_FOUND") || strings.Contains(message, "TRAINER_NOT_FOUND") {
			statusCode = fiber.StatusNotFound
		} else if strings.Contains(message, "INVALID_PRODUCT") {
			statusCode = fiber.StatusBadRequest
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"status":      "error",
			"status_code": statusCode,
			"message":     message,
			"result":      nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "success",
		"status_code": fiber.StatusCreated,
		"message":     "Session package renewed successfully",
		"result":      result,
	})
}
