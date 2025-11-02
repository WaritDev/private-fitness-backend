package controller

import (
	"strings"
	"time"
	"fmt"
	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/services"
	"github.com/gofiber/fiber/v2"
)

type TrainerController struct {
	sessionUC *services.SessionService
	trainerUC *services.TrainerService
}

func ProvideTrainerController(sessionUC *services.SessionService, trainerUC *services.TrainerService) *TrainerController {
	return &TrainerController{
		sessionUC: sessionUC,
		trainerUC: trainerUC,
	}
}

// POST /api/trainers/match - จับคู่เทรนเนอร์ตาม Use Case 4S
func (h *TrainerController) MatchTrainer(c *fiber.Ctx) error {
	// Parse raw JSON first to handle time parsing manually if needed
	var rawReq map[string]interface{}
	if err := c.BodyParser(&rawReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid request body: " + err.Error(),
			"result":      nil,
		})
	}

	var req requests.MatchTrainerRequest

	// Parse dayOfWeek
	if dayOfWeek, ok := rawReq["dayOfWeek"].(string); ok {
		req.DayOfWeek = dayOfWeek
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "dayOfWeek is required and must be a string",
			"result":      nil,
		})
	}

	// Parse startTime from string
	if startTimeStr, ok := rawReq["startTime"].(string); ok {
		startTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			// Try alternative ISO 8601 format
			startTime, err = time.Parse("2006-01-02T15:04:05Z07:00", startTimeStr)
			if err != nil {
				startTime, err = time.Parse("2006-01-02T15:04:05.000Z", startTimeStr)
				if err != nil {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"status":      "Bad Request",
						"status_code": fiber.StatusBadRequest,
						"message":     "Invalid startTime format. Expected ISO 8601 format. Error: " + err.Error(),
						"result":      nil,
					})
				}
			}
		}
		req.StartTime = startTime
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "startTime is required and must be a string (ISO 8601 format)",
			"result":      nil,
		})
	}

	// Parse endTime from string
	if endTimeStr, ok := rawReq["endTime"].(string); ok {
		endTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			// Try alternative ISO 8601 format
			endTime, err = time.Parse("2006-01-02T15:04:05Z07:00", endTimeStr)
			if err != nil {
				endTime, err = time.Parse("2006-01-02T15:04:05.000Z", endTimeStr)
				if err != nil {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"status":      "Bad Request",
						"status_code": fiber.StatusBadRequest,
						"message":     "Invalid endTime format. Expected ISO 8601 format. Error: " + err.Error(),
						"result":      nil,
					})
				}
			}
		}
		req.EndTime = endTime
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "Bad Request",
			"status_code": fiber.StatusBadRequest,
			"message":     "endTime is required and must be a string (ISO 8601 format)",
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
			"message":     "startTime and endTime are required and must be valid",
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
	fmt.Println("result", result)
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
func (h *TrainerController) ListAllTrainers(c *fiber.Ctx) error {
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

// Use Case 1P: Manage Working Hours

// GET /api/trainers/working-hours - ดึงรายการเวลาทำงานของเทรนเนอร์
func (h *TrainerController) GetWorkingHours(c *fiber.Ctx) error {
	// ดึง trainerUsername จาก JWT claims (set by middleware)
	trainerUsername, ok := c.Locals("username").(string)
	if !ok || trainerUsername == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized - JWT token required",
		})
	}

	// เรียก use case
	result, err := h.trainerUC.GetWorkingHours(c.Context(), trainerUsername)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// POST /api/trainers/working-hours - เพิ่มเวลาทำงานใหม่
func (h *TrainerController) AddWorkingTime(c *fiber.Ctx) error {
	// ดึง trainerUsername จาก JWT claims (set by middleware)
	trainerUsername, ok := c.Locals("username").(string)
	if !ok || trainerUsername == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized - JWT token required",
		})
	}

	// Parse request body
	var req requests.AddWorkingTimeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// Validate required fields (Fiber validator will check this)
	if req.DayOfWeek == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "dayOfWeek is required",
		})
	}

	if req.StartTime == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "startTime is required",
		})
	}

	if req.EndTime == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "endTime is required",
		})
	}

	// Validate day of week
	validDays := []string{"MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY", "SATURDAY", "SUNDAY"}
	isValidDay := false
	for _, day := range validDays {
		if req.DayOfWeek == day {
			isValidDay = true
			break
		}
	}
	if !isValidDay {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid dayOfWeek. Must be MONDAY, TUESDAY, WEDNESDAY, THURSDAY, FRIDAY, SATURDAY, or SUNDAY",
		})
	}

	// เรียก use case (จะทำ validation ทั้งหมดใน use case)
	result, err := h.trainerUC.AddWorkingTime(c.Context(), trainerUsername, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// ถ้า result status = error แสดงว่าเกิด validation error
	if result.Status == "error" {
		return c.Status(fiber.StatusBadRequest).JSON(result)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// PUT /api/trainers/working-hours/:id - แก้ไขเวลาทำงาน (Q1P.4)
func (h *TrainerController) UpdateWorkingTime(c *fiber.Ctx) error {
	// ดึง trainerUsername จาก JWT claims (set by middleware)
	trainerUsername, ok := c.Locals("username").(string)
	if !ok || trainerUsername == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized - JWT token required",
		})
	}

	// ดึง ID จาก URL parameter
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid ID parameter",
		})
	}

	// Parse request body
	var req requests.UpdateWorkingTimeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// Validate required fields
	if req.DayOfWeek == "" || req.StartTime == "" || req.EndTime == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "dayOfWeek, startTime, and endTime are required",
		})
	}

	// Validate day of week
	validDays := []string{"MONDAY", "TUESDAY", "WEDNESDAY", "THURSDAY", "FRIDAY", "SATURDAY", "SUNDAY"}
	isValidDay := false
	for _, day := range validDays {
		if req.DayOfWeek == day {
			isValidDay = true
			break
		}
	}
	if !isValidDay {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid dayOfWeek. Must be MONDAY, TUESDAY, WEDNESDAY, THURSDAY, FRIDAY, SATURDAY, or SUNDAY",
		})
	}

	// เรียก use case
	result, err := h.trainerUC.UpdateWorkingTime(c.Context(), trainerUsername, int32(id), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// ถ้า result status = error แสดงว่าเกิด validation error
	if result.Status == "error" {
		return c.Status(fiber.StatusBadRequest).JSON(result)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// DELETE /api/trainers/working-hours/:id - ลบเวลาทำงาน (Q1P.5)
func (h *TrainerController) DeleteWorkingTime(c *fiber.Ctx) error {
	// ดึง trainerUsername จาก JWT claims (set by middleware)
	trainerUsername, ok := c.Locals("username").(string)
	if !ok || trainerUsername == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized - JWT token required",
		})
	}

	// ดึง ID จาก URL parameter
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid ID parameter",
		})
	}

	// เรียก use case
	result, err := h.trainerUC.DeleteWorkingTime(c.Context(), trainerUsername, int32(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// ถ้า result status = error แสดงว่าเกิด validation error
	if result.Status == "error" {
		return c.Status(fiber.StatusNotFound).JSON(result)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// ========== Use Case: Trainer Calendar & Check-in Confirmation ==========

// GET /api/trainers/calendar - ดึง appointments พร้อม pending check-ins
func (h *TrainerController) GetCalendar(c *fiber.Ctx) error {
	// ดึง trainerUsername จาก JWT claims (set by middleware)
	trainerUsername, ok := c.Locals("username").(string)
	if !ok || trainerUsername == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized - JWT token required",
		})
	}

	// เรียก use case
	result, err := h.trainerUC.GetCalendar(c.Context(), trainerUsername)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// POST /api/trainers/checkin - Trainer confirm check-in และหัก session
func (h *TrainerController) ConfirmCheckIn(c *fiber.Ctx) error {
	// ดึง trainerUsername จาก JWT claims (set by middleware)
	trainerUsername, ok := c.Locals("username").(string)
	if !ok || trainerUsername == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized - JWT token required",
		})
	}

	// Parse request body
	var req requests.ConfirmCheckInRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// Validate required fields
	if req.SessionID == 0 || req.CustomerUsername == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "sessionId and customerUsername are required",
		})
	}

	// เรียก use case
	result, err := h.trainerUC.ConfirmCheckIn(c.Context(), trainerUsername, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// ถ้า result status = error แสดงว่าเกิด validation error
	if result.Status == "error" {
		return c.Status(fiber.StatusBadRequest).JSON(result)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// ========== Use Case 3P: Manage Day-Offs ==========

// GET /api/trainers/day-offs - ดึงรายการวันหยุดของเทรนเนอร์
func (h *TrainerController) GetDayOffs(c *fiber.Ctx) error {
	// ดึง trainerUsername จาก JWT claims (set by middleware)
	trainerUsername, ok := c.Locals("username").(string)
	if !ok || trainerUsername == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized - JWT token required",
		})
	}

	// เรียก use case
	result, err := h.trainerUC.GetDayOffs(c.Context(), trainerUsername)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// POST /api/trainers/day-offs - เพิ่มวันหยุดใหม่
func (h *TrainerController) AddDayOff(c *fiber.Ctx) error {
	// ดึง trainerUsername จาก JWT claims (set by middleware)
	trainerUsername, ok := c.Locals("username").(string)
	if !ok || trainerUsername == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized - JWT token required",
		})
	}

	// Parse request body
	var req requests.AddDayOffRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// Validate required field
	if req.DayOffDate == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "dayOffDate is required",
		})
	}

	// เรียก use case (จะทำ validation ทั้งหมดใน use case)
	result, err := h.trainerUC.AddDayOff(c.Context(), trainerUsername, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// ถ้า result status = error แสดงว่าเกิด validation error
	if result.Status == "error" {
		return c.Status(fiber.StatusBadRequest).JSON(result)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// DELETE /api/trainers/day-offs/:id - ลบวันหยุด
func (h *TrainerController) DeleteDayOff(c *fiber.Ctx) error {
	// ดึง trainerUsername จาก JWT claims (set by middleware)
	trainerUsername, ok := c.Locals("username").(string)
	if !ok || trainerUsername == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized - JWT token required",
		})
	}

	// ดึง ID จาก URL parameter
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid ID parameter",
		})
	}

	// เรียก use case
	result, err := h.trainerUC.DeleteDayOff(c.Context(), trainerUsername, int32(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// ถ้า result status = error แสดงว่าเกิด validation error
	if result.Status == "error" {
		return c.Status(fiber.StatusNotFound).JSON(result)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
