package rest

import (
	"strings"

	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/usecases"
	"github.com/gofiber/fiber/v2"
)

type TrainerHandler struct {
	sessionUC *usecases.SessionUseCase
}

func ProvideTrainerHandler(sessionUC *usecases.SessionUseCase) *TrainerHandler {
	return &TrainerHandler{
		sessionUC: sessionUC,
	}
}

// POST /api/trainers/match - จับคู่เทรนเนอร์ตาม Use Case 4S
func (h *TrainerHandler) MatchTrainer(c *fiber.Ctx) error {
	var req requests.MatchTrainerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid request body",
			"result":      nil,
		})
	}

	// Validate day of week
	validDays := []string{"MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY", "SATURDAY", "SUNDAY"}
	isValid := false
	for _, day := range validDays {
		if req.DayOfWeek == day {
			isValid = true
			break
		}
	}
	if !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid dayOfWeek. Must be one of: MONDAY, TUESDAY, WEDNESDAY, THURSDAY, FRIDAY, SATURDAY, SUNDAY",
			"result":      nil,
		})
	}

	// Validate times
	if req.StartTime.IsZero() || req.EndTime.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "startTime and endTime are required",
			"result":      nil,
		})
	}

	if req.EndTime.Before(req.StartTime) || req.EndTime.Equal(req.StartTime) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "endTime must be after startTime",
			"result":      nil,
		})
	}

	// Call use case
	result, err := h.sessionUC.MatchTrainer(c.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "NO_TRAINER_AVAILABLE") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":      "Not Found",
				"status_code": fiber.StatusNotFound,
				"message":     "No available trainer found for the selected day and time",
				"result":      nil,
			})
		}

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
		"message":     "Trainer matched successfully",
		"result":      result,
	})
}

// GET /api/trainers - ดึงรายชื่อเทรนเนอร์ทั้งหมด (สำหรับ dropdown)
func (h *TrainerHandler) ListAllTrainers(c *fiber.Ctx) error {
	result, err := h.sessionUC.ListAllTrainers(c.Context())
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
		"message":     "Trainers retrieved successfully",
		"result":      result,
	})
}
