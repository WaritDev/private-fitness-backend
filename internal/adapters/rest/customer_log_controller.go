package rest

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/usecases"
)

type CustomerLogHandler struct {
	uc *usecases.CustomerLogUsecase
}

func ProvideCustomerLogHandler(uc *usecases.CustomerLogUsecase) *CustomerLogHandler {
	return &CustomerLogHandler{uc: uc}
}

// GET /api/customer-logs?page=1&limit=10
func (h *CustomerLogHandler) List(c *fiber.Ctx) error {
	res, err := h.uc.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// POST /api/customer-logs/:id/update
func (h *CustomerLogHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var req requests.UpdateCustomerLogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}

	res, err := h.uc.Update(c.Context(), int32(id64), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *CustomerLogHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	res, err := h.uc.Delete(c.Context(), int32(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// GET /api/customer-logs/:id
func (h *CustomerLogHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id", "")
	resp, err := h.uc.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}