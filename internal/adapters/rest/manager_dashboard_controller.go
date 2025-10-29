package rest

import (
	"time"

	req "github.com/WaritDev/private-fitness-backend/domain/requests"
	uc "github.com/WaritDev/private-fitness-backend/domain/usecases"
	"github.com/gofiber/fiber/v2"
)

type ManagerDashboardHandler struct {
	uc uc.ManagerDashboardUsecase
}

func ProvideManagerDashboardHandler(u uc.ManagerDashboardUsecase) *ManagerDashboardHandler {
	return &ManagerDashboardHandler{uc: u}
}

// GET/POST /api/manager/dashboard?start=YYYY-MM-DD&end=YYYY-MM-DD
func (h *ManagerDashboardHandler) GetDashboard(c *fiber.Ctx) error {
	var (
		start time.Time
		end   time.Time
	)

	if s := c.Query("start"); s != "" {
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid start"})
		}
		start = t
	}
	if e := c.Query("end"); e != "" {
		t, err := time.Parse("2006-01-02", e)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid end"})
		}
		end = t
	}

	if start.IsZero() || end.IsZero() {
		end = time.Now()
		start = end.AddDate(0, 0, -30)
	}

	if c.Method() == fiber.MethodPost && len(c.Body()) > 0 {
		var body req.ManagerDashboardRequest
		if err := c.BodyParser(&body); err == nil {
			if body.Start != nil {
				start = *body.Start
			}
			if body.End != nil {
				end = *body.End
			}
		}
	}

	resp, err := h.uc.Get(c.Context(), start, end)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}