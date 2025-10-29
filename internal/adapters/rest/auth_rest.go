package rest

import (
	"github.com/WaritDev/private-fitness-backend/domain/usecases"
	"github.com/gofiber/fiber/v2"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

type AuthHandler struct {
	UC *usecases.AuthUseCase
}

func ProvideAuthHandler(uc *usecases.AuthUseCase) *AuthHandler {
	return &AuthHandler{UC: uc}
}

// POST /api/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req := new(dbmodel.UsersCredentials)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"status":      fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}
	user, err := h.UC.Login(c.Context(), req)
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     err.Error(),
			"result":      nil,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(user)
}