package controller

import (
	"time"

	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/services"
	"github.com/gofiber/fiber/v2"
)

type BookingController struct {
	useCase *services.BookingService
}

func ProvideBookingController(useCase *services.BookingService) *BookingController {
	return &BookingController{
		useCase: useCase,
	}
}

// GetBookingSlots - GET /api/bookings/slots
// ดึงข้อมูล Booking Slots พร้อมข้อมูล availability, day offs, และ appointments
func (h *BookingController) GetBookingSlots(c *fiber.Ctx) error {
	// Parse query parameters
	trainerUsername := c.Query("trainerUsername")
	calendarStartStr := c.Query("calendarStart") // ISO format: 2025-11-01T00:00:00Z
	calendarEndStr := c.Query("calendarEnd")
	customerUsername := c.Query("customerUsername") // Optional

	// Validate required fields
	if trainerUsername == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusBadRequest,
			"message":     "trainerUsername is required",
			"result":      nil,
		})
	}

	if calendarStartStr == "" || calendarEndStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusBadRequest,
			"message":     "calendarStart and calendarEnd are required",
			"result":      nil,
		})
	}

	// Parse dates
	calendarStart, err := time.Parse(time.RFC3339, calendarStartStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid calendarStart format, expected ISO 8601 (e.g., 2025-11-01T00:00:00Z)",
			"result":      nil,
		})
	}

	calendarEnd, err := time.Parse(time.RFC3339, calendarEndStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid calendarEnd format, expected ISO 8601 (e.g., 2025-11-30T23:59:59Z)",
			"result":      nil,
		})
	}

	// Build request
	req := requests.GetBookingSlotsRequest{
		TrainerUsername:  trainerUsername,
		CalendarStart:    calendarStart,
		CalendarEnd:      calendarEnd,
		CustomerUsername: customerUsername,
	}

	// Call use case
	result, err := h.useCase.GetBookingSlots(c.Context(), req)
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
		"message":     "Booking slots retrieved successfully",
		"result":      result,
	})
}

// BookAppointment - POST /api/bookings/book
// Q3C: จองนัดหมายกับเทรนเนอร์ (Transaction: Check + INSERT + UPDATE + LOG)
func (h *BookingController) BookAppointment(c *fiber.Ctx) error {
	// Parse request body
	var req requests.BookAppointmentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid request body",
			"result":      nil,
		})
	}

	// Validate required fields (sessionId is now optional - will auto-find active session)
	if req.TrainerUsername == "" || req.CustomerUsername == "" || req.StartTime == "" || req.EndTime == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Required fields: trainerUsername, customerUsername, startTime, endTime",
			"result":      nil,
		})
	}

	// Call use case
	result, err := h.useCase.BookAppointment(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusInternalServerError,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	// Check if booking was successful
	if !result.Success {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusBadRequest,
			"message":     result.Message,
			"result":      result,
		})
	}

	// Success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "success",
		"status_code": fiber.StatusCreated,
		"message":     result.Message,
		"result":      result,
	})
}

// CancelAppointment - DELETE /api/bookings/cancel/:id
// ยกเลิกการจอง (Transaction: Check + DELETE + DECREMENT + LOG)
func (h *BookingController) CancelAppointment(c *fiber.Ctx) error {
	// Parse path parameter
	appointmentID, err := c.ParamsInt("id")
	if err != nil || appointmentID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid appointment ID",
			"result":      nil,
		})
	}

	// Parse request body (need customerUsername for authorization)
	var req requests.CancelAppointmentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusBadRequest,
			"message":     "Invalid request body",
			"result":      nil,
		})
	}

	// Set appointmentID from path parameter
	req.AppointmentID = int32(appointmentID)

	// Validate required fields
	if req.CustomerUsername == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusBadRequest,
			"message":     "customerUsername is required",
			"result":      nil,
		})
	}

	// Call use case
	result, err := h.useCase.CancelAppointment(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusInternalServerError,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	// Check if cancellation was successful
	if !result.Success {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":      "error",
			"status_code": fiber.StatusBadRequest,
			"message":     result.Message,
			"result":      result,
		})
	}

	// Success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "success",
		"status_code": fiber.StatusOK,
		"message":     result.Message,
		"result":      result,
	})
}
