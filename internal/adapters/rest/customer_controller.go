package rest

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/usecases"
)

type CustomerHandler struct {
	uc *usecases.CustomerUsecase
}

func ProvideCustomerHandler(uc *usecases.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{uc: uc}
}

func (h *CustomerHandler) ListCustomers(c *fiber.Ctx) error {
	resp, err := h.uc.ListCustomers(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *CustomerHandler) UpdateCustomer(c *fiber.Ctx) error {
	username := c.Params("username")
	var body requests.UpdateCustomerRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if body.ConfirmNewPassword == nil {
		if v := c.Query("confirmPassword"); v != "" { body.ConfirmNewPassword = &[]string{v}[0] }
	}

	resp, err := h.uc.UpdateCustomer(c.Context(), username, body)
	if err != nil { return fiber.NewError(fiber.StatusBadRequest, err.Error()) }
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *CustomerHandler) DeleteCustomer(c *fiber.Ctx) error {
	username := c.Params("username")
	resp, err := h.uc.DeleteCustomer(c.Context(), username)
	if err != nil { return fiber.NewError(fiber.StatusBadRequest, err.Error()) }
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *CustomerHandler) GetByUsername(c *fiber.Ctx) error {
	username := strings.TrimSpace(c.Params("username"))
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "username required"})
	}

	out, err := h.uc.GetCustomerByUsername(c.Context(), username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "customer not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.JSON(out)
}