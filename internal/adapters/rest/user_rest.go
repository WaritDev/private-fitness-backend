package rest

import (
	"github.com/WaritDev/private-fitness-backend/domain/usecases"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	UC *usecases.UserUseCase
}

func ProvideUserHandler(uc *usecases.UserUseCase) *UserHandler {
	return &UserHandler{UC: uc}
}

// GET /api/users/check-phone?phone=0812345678
// ตรวจสอบว่าเบอร์โทรศัพท์ซ้ำหรือไม่ (Q3S.1)
func (h *UserHandler) CheckPhoneNumber(c *fiber.Ctx) error {
	phoneNumber := c.Query("phone")
	if phoneNumber == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "Phone number is required",
			"result":      nil,
		})
	}

	exists, err := h.UC.CheckPhoneNumberExistsUser(c.Context(), phoneNumber)
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
		"message":     "Phone number check completed",
		"result": fiber.Map{
			"exists":      exists,
			"available":   !exists,
			"phoneNumber": phoneNumber,
		},
	})
}

// GET /api/users/check-gmail?gmail=test@example.com
// ตรวจสอบว่าอีเมลซ้ำหรือไม่ (Q3S.2)
func (h *UserHandler) CheckGmail(c *fiber.Ctx) error {
	gmail := c.Query("gmail")
	if gmail == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "Gmail is required",
			"result":      nil,
		})
	}

	exists, err := h.UC.CheckGmailExistsUser(c.Context(), gmail)
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
		"message":     "Gmail check completed",
		"result": fiber.Map{
			"exists":    exists,
			"available": !exists,
			"gmail":     gmail,
		},
	})
}

// GET /api/users/check-username?username=testuser
// ตรวจสอบว่าชื่อผู้ใช้งานซ้ำหรือไม่ (Q3S.3)
func (h *UserHandler) CheckUsername(c *fiber.Ctx) error {
	username := c.Query("username")
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "Username is required",
			"result":      nil,
		})
	}

	exists, err := h.UC.CheckUsernameExists(c.Context(), username)
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
		"message":     "Username check completed",
		"result": fiber.Map{
			"exists":    exists,
			"available": !exists,
			"username":  username,
		},
	})
}