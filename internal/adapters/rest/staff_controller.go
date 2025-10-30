package rest

import (
	"github.com/gofiber/fiber/v2"

	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/usecases"
)

type StaffHandler struct {
	uc *usecases.StaffUsecase
}

func ProvideStaffHandler(uc *usecases.StaffUsecase) *StaffHandler {
	return &StaffHandler{uc: uc}
}

func (h *StaffHandler) Register(r fiber.Router) {
	r.Get("/staffs", h.ListStaffs)
	r.Post("/staffs", h.CreateStaff)
	r.Post("/staffs/:username/update", h.UpdateStaff)
	r.Delete("/staffs/:username", h.DeleteStaff)
}

func (h *StaffHandler) ListStaffs(c *fiber.Ctx) error {
	var q requests.ListStaffsRequest
	if err := c.QueryParser(&q); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	resp, err := h.uc.ListStaffs(c.Context(), q)
	if err != nil { return fiber.NewError(fiber.StatusBadRequest, err.Error()) }
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *StaffHandler) CreateStaff(c *fiber.Ctx) error {
	var body requests.CreateStaffRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	resp, err := h.uc.CreateStaff(c.Context(), body)
	if err != nil { return fiber.NewError(fiber.StatusBadRequest, err.Error()) }
	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *StaffHandler) UpdateStaff(c *fiber.Ctx) error {
	target := c.Params("username")
	var body requests.UpdateStaffRequest
	if err := c.BodyParser(&body); err != nil { return fiber.NewError(fiber.StatusBadRequest, err.Error()) }

	resp, err := h.uc.UpdateStaff(c.Context(), target, body)
	if err != nil { return fiber.NewError(fiber.StatusBadRequest, err.Error()) }
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *StaffHandler) DeleteStaff(c *fiber.Ctx) error {
	target := c.Params("username")
	resp, err := h.uc.DeleteStaff(c.Context(), target)
	if err != nil { return fiber.NewError(fiber.StatusBadRequest, err.Error()) }
	return c.Status(fiber.StatusOK).JSON(resp)
}