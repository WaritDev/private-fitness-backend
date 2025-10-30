package responses

import "time"

// BookingSlot - Individual booking slot (30-minute interval)
type BookingSlot struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Available bool      `json:"available"`
	IsBooked  bool      `json:"isBooked"` // true if this customer booked this slot
	BookedBy  string    `json:"bookedBy"` // Customer username (only show if it's this customer)
	SlotType  string    `json:"slotType"` // "available", "booked", "day_off", "unavailable"
}

// DayOffSlot - Day off or unavailable time slot
type DayOffSlot struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

// AppointmentSlot - Booked appointment slot
type AppointmentSlot struct {
	StartTime        time.Time `json:"startTime"`
	EndTime          time.Time `json:"endTime"`
	CustomerUsername string    `json:"customerUsername"`
}

// TrainerAvailability - Trainer's weekly working hours
type TrainerAvailability struct {
	DayOfWeek string    `json:"dayOfWeek"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

// GetBookingSlotsResponse - Response with all booking information
type GetBookingSlotsResponse struct {
	TrainerUsername    string                `json:"trainerUsername"`
	CalendarStart      time.Time             `json:"calendarStart"`
	CalendarEnd        time.Time             `json:"calendarEnd"`
	WeeklyAvailability []TrainerAvailability `json:"weeklyAvailability"`
	DayOffSlots        []DayOffSlot          `json:"dayOffSlots"`
	BookedAppointments []AppointmentSlot     `json:"bookedAppointments"`
	AvailableSlots     []BookingSlot         `json:"availableSlots"`   // All available 30-min slots
	CustomerBookings   []BookingSlot         `json:"customerBookings"` // This customer's bookings only
	Message            string                `json:"message"`
}

// BookAppointmentResponse - Q3C: Response after booking appointment
type BookAppointmentResponse struct {
	Success          bool      `json:"success"`
	Message          string    `json:"message"`
	AppointmentID    int32     `json:"appointmentId,omitempty"`
	TrainerUsername  string    `json:"trainerUsername"`
	CustomerUsername string    `json:"customerUsername"`
	StartTime        time.Time `json:"startTime"`
	EndTime          time.Time `json:"endTime"`
	SessionID        int32     `json:"sessionId"`
	RemainingSession int32     `json:"remainingSessions,omitempty"` // จำนวนครั้งที่เหลือหลังจากจอง
}

// CancelAppointmentResponse - Response after canceling appointment
type CancelAppointmentResponse struct {
	Success          bool      `json:"success"`
	Message          string    `json:"message"`
	AppointmentID    int32     `json:"appointmentId"`
	CustomerUsername string    `json:"customerUsername"`
	StartTime        time.Time `json:"startTime,omitempty"`
	EndTime          time.Time `json:"endTime,omitempty"`
	SessionID        int32     `json:"sessionId,omitempty"`
	RemainingSession int32     `json:"remainingSessions,omitempty"` // จำนวนครั้งที่เหลือหลังจากยกเลิก
}
