package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/services"
)

type CustomerLogController struct {
	uc *services.CustomerLogService
}

func ProvideCustomerLogController(uc *services.CustomerLogService) *CustomerLogController {
	return &CustomerLogController{uc: uc}
}

// GET /api/customer-logs?page=1&limit=10
func (h *CustomerLogController) List(c *fiber.Ctx) error {
	res, err := h.uc.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// POST /api/customer-logs/:id/update
func (h *CustomerLogController) Update(c *fiber.Ctx) error {
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

func (h *CustomerLogController) Delete(c *fiber.Ctx) error {
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
func (h *CustomerLogController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id", "")
	resp, err := h.uc.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}