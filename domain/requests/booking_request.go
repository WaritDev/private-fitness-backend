package requests

import "time"

// GetBookingSlotsRequest - Request for getting available booking slots
type GetBookingSlotsRequest struct {
	TrainerUsername  string    `json:"trainerUsername" validate:"required"`
	CalendarStart    time.Time `json:"calendarStart" validate:"required"`
	CalendarEnd      time.Time `json:"calendarEnd" validate:"required"`
	CustomerUsername string    `json:"customerUsername"` // Optional: to show customer's own bookings
}

// BookAppointmentRequest - Q3C: Request for booking an appointment
type BookAppointmentRequest struct {
	TrainerUsername  string `json:"trainerUsername" validate:"required"`
	CustomerUsername string `json:"customerUsername" validate:"required"`
	SessionID        *int32 `json:"sessionId,omitempty"`           // Optional: ถ้าไม่ระบุจะหา active session อัตโนมัติ
	StartTime        string `json:"startTime" validate:"required"` // ISO 8601 format
	EndTime          string `json:"endTime" validate:"required"`   // ISO 8601 format
}

// CancelAppointmentRequest - Request for canceling an appointment
type CancelAppointmentRequest struct {
	AppointmentID    int32  `json:"appointmentId" validate:"required,min=1"`
	CustomerUsername string `json:"customerUsername" validate:"required"` // เพื่อตรวจสอบว่าเป็นเจ้าของการจองจริง
}
