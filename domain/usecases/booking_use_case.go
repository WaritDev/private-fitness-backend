package usecases

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/responses"
)

type BookingUseCase struct {
	scheduleRepo    repositories.TrainingScheduleRepository
	sessionRepo     repositories.CustomerSessionRepository
	customerLogRepo repositories.CustomerLogRepository
	db              *sql.DB
}

func ProvideBookingUseCase(
	scheduleRepo repositories.TrainingScheduleRepository,
	sessionRepo repositories.CustomerSessionRepository,
	customerLogRepo repositories.CustomerLogRepository,
	db *sql.DB,
) *BookingUseCase {
	return &BookingUseCase{
		scheduleRepo:    scheduleRepo,
		sessionRepo:     sessionRepo,
		customerLogRepo: customerLogRepo,
		db:              db,
	}
}

// GetBookingSlots - ดึงข้อมูล Booking Slots โดยรวมข้อมูล 3 ส่วน
func (u *BookingUseCase) GetBookingSlots(ctx context.Context, req requests.GetBookingSlotsRequest) (*responses.GetBookingSlotsResponse, error) {
	// 1. ดึงเวลาทำงานประจำสัปดาห์ (TRAINING_AVAILABILITY)
	availability, err := u.scheduleRepo.GetTrainerAvailability(ctx, req.TrainerUsername)
	if err != nil {
		return nil, fmt.Errorf("failed to get trainer availability: %w", err)
	}

	// 2. ดึงวันหยุด/ช่วงเวลาที่ไม่รับนัด (DAY_OFF)
	dayOffs, err := u.scheduleRepo.GetDayOffSchedules(ctx, req.TrainerUsername, req.CalendarStart, req.CalendarEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to get day off schedules: %w", err)
	}

	// 3. ดึงนัดที่ถูกจองแล้ว (APPOINTMENT)
	appointments, err := u.scheduleRepo.GetAppointmentSchedules(ctx, req.TrainerUsername, req.CalendarStart, req.CalendarEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment schedules: %w", err)
	}

	// 4. คำนวณ Available Slots (30-minute intervals)
	availableSlots := u.calculateAvailableSlots(availability, dayOffs, appointments, req.CalendarStart, req.CalendarEnd, req.CustomerUsername)

	// 5. Filter customer's own bookings
	customerBookings := u.filterCustomerBookings(appointments, req.CustomerUsername)

	// Build response
	return &responses.GetBookingSlotsResponse{
		TrainerUsername:    req.TrainerUsername,
		CalendarStart:      req.CalendarStart,
		CalendarEnd:        req.CalendarEnd,
		WeeklyAvailability: u.mapToWeeklyAvailability(availability),
		DayOffSlots:        u.mapToDayOffSlots(dayOffs),
		BookedAppointments: u.mapToAppointmentSlots(appointments),
		AvailableSlots:     availableSlots,
		CustomerBookings:   customerBookings,
		Message:            "Booking slots retrieved successfully",
	}, nil
}

// calculateAvailableSlots - คำนวณ Booking Slots ทุก 30 นาที
func (u *BookingUseCase) calculateAvailableSlots(
	availability []repositories.TrainerAvailabilityInfo,
	dayOffs []repositories.ScheduleTimeSlot,
	appointments []repositories.AppointmentInfo,
	calendarStart, calendarEnd time.Time,
	customerUsername string,
) []responses.BookingSlot {
	slots := []responses.BookingSlot{}

	// TODO: Implement slot calculation logic
	// 1. นำเวลาทำงานจาก availability มาสร้าง slots ทุก 30 นาที
	// 2. ลบ slots ที่ตรงกับ dayOffs
	// 3. ทำเครื่องหมาย slots ที่ถูกจองแล้วใน appointments
	// 4. แสดงเฉพาะ slots ของ customer นี้เอง (ถ้ามี customerUsername)

	return slots
}

// filterCustomerBookings - กรองเฉพาะการจองของ customer นี้
func (u *BookingUseCase) filterCustomerBookings(
	appointments []repositories.AppointmentInfo,
	customerUsername string,
) []responses.BookingSlot {
	if customerUsername == "" {
		return []responses.BookingSlot{}
	}

	bookings := []responses.BookingSlot{}
	for _, apt := range appointments {
		if apt.CustomerUsername == customerUsername {
			bookings = append(bookings, responses.BookingSlot{
				StartTime: apt.StartTime,
				EndTime:   apt.EndTime,
				Available: false,
				IsBooked:  true,
				BookedBy:  customerUsername,
				SlotType:  "booked",
			})
		}
	}

	return bookings
}

// mapToWeeklyAvailability - แปลงข้อมูล availability เป็น response format
func (u *BookingUseCase) mapToWeeklyAvailability(availability []repositories.TrainerAvailabilityInfo) []responses.TrainerAvailability {
	result := make([]responses.TrainerAvailability, len(availability))
	for i, a := range availability {
		result[i] = responses.TrainerAvailability{
			DayOfWeek: a.DayOfWeek,
			StartTime: a.StartTime,
			EndTime:   a.EndTime,
		}
	}
	return result
}

// mapToDayOffSlots - แปลงข้อมูล day offs เป็น response format
func (u *BookingUseCase) mapToDayOffSlots(dayOffs []repositories.ScheduleTimeSlot) []responses.DayOffSlot {
	result := make([]responses.DayOffSlot, len(dayOffs))
	for i, d := range dayOffs {
		result[i] = responses.DayOffSlot{
			StartTime: d.StartTime,
			EndTime:   d.EndTime,
		}
	}
	return result
}

// mapToAppointmentSlots - แปลงข้อมูล appointments เป็น response format
func (u *BookingUseCase) mapToAppointmentSlots(appointments []repositories.AppointmentInfo) []responses.AppointmentSlot {
	result := make([]responses.AppointmentSlot, len(appointments))
	for i, a := range appointments {
		result[i] = responses.AppointmentSlot{
			StartTime:        a.StartTime,
			EndTime:          a.EndTime,
			CustomerUsername: a.CustomerUsername,
		}
	}
	return result
}

// BookAppointment - Q3C: จองนัดหมาย (Transaction: Check + INSERT schedule + UPDATE session + INSERT log)
func (u *BookingUseCase) BookAppointment(ctx context.Context, req requests.BookAppointmentRequest) (*responses.BookAppointmentResponse, error) {
	// Parse time strings to time.Time
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start time format: %w", err)
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end time format: %w", err)
	}

	// Validation 1: หา Session package ACTIVE ของ Customer
	activeSession, err := u.sessionRepo.GetActiveSessionByCustomer(ctx, req.CustomerUsername)
	if err != nil {
		return nil, fmt.Errorf("failed to get active session: %w", err)
	}
	if activeSession == nil {
		return &responses.BookAppointmentResponse{
			Success: false,
			Message: "Customer does not have an active session package or no sessions remaining",
		}, nil
	}

	// ใช้ session_id ที่ระบุมา หรือใช้ session_id ที่หาได้
	sessionID := activeSession.ID
	if req.SessionID != nil {
		sessionID = *req.SessionID
	}

	// Validation 2: Q3C.5 - ตรวจสอบว่าช่วงเวลายังว่างอยู่จริง
	overlappedCount, err := u.scheduleRepo.CheckTimeSlotAvailability(ctx, req.TrainerUsername, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to check time slot availability: %w", err)
	}
	if overlappedCount > 0 {
		return &responses.BookAppointmentResponse{
			Success: false,
			Message: fmt.Sprintf("Time slot is not available. Found %d overlapping appointment(s)", overlappedCount),
		}, nil
	}

	// Q3C.6 - Execute Transaction: BookAppointment + IncrementUsedSessions + CreateCustomerLog
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Insert TRAINING_SCHEDULE (APPOINTMENT)
	err = u.scheduleRepo.BookAppointment(ctx, repositories.BookAppointmentParams{
		TrainerUsername:  req.TrainerUsername,
		CustomerUsername: req.CustomerUsername,
		SessionID:        sessionID,
		StartTime:        startTime,
		EndTime:          endTime,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to book appointment: %w", err)
	}

	// 2. Update CUSTOMER_SESSION (used_sessions + 1)
	err = u.sessionRepo.IncrementUsedSessions(ctx, tx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to increment used sessions: %w", err)
	}

	// 3. Insert CUSTOMER_LOG (log_type = 'BOOK_SESSION')
	err = u.customerLogRepo.CreateCustomerLog(ctx, tx, req.CustomerUsername, "BOOK_SESSION")
	if err != nil {
		return nil, fmt.Errorf("failed to create customer log: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// คำนวณ remaining sessions
	remainingSessions := activeSession.TotalSessions - (activeSession.UsedSessions + 1)

	// Return success response
	return &responses.BookAppointmentResponse{
		Success:          true,
		Message:          "Appointment booked successfully",
		TrainerUsername:  req.TrainerUsername,
		CustomerUsername: req.CustomerUsername,
		StartTime:        startTime,
		EndTime:          endTime,
		SessionID:        sessionID,
		RemainingSession: remainingSessions,
	}, nil
}

// CancelAppointment - ยกเลิกการจอง (Transaction: Check + DELETE schedule + DECREMENT session + INSERT log)
func (u *BookingUseCase) CancelAppointment(ctx context.Context, req requests.CancelAppointmentRequest) (*responses.CancelAppointmentResponse, error) {
	// Validation 1: ตรวจสอบว่าการจองนี้มีอยู่จริง
	appointment, err := u.scheduleRepo.GetAppointmentById(ctx, req.AppointmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}
	if appointment == nil {
		return &responses.CancelAppointmentResponse{
			Success: false,
			Message: "Appointment not found",
		}, nil
	}

	// Validation 2: ตรวจสอบว่าเป็นเจ้าของการจองจริง
	if appointment.CustomerUsername != req.CustomerUsername {
		return &responses.CancelAppointmentResponse{
			Success: false,
			Message: "You are not authorized to cancel this appointment",
		}, nil
	}

	// Validation 3: ตรวจสอบว่าการจองยังไม่ผ่านไปแล้ว (optional - ขึ้นอยู่กับ business rule)
	if time.Now().After(appointment.StartTime) {
		return &responses.CancelAppointmentResponse{
			Success: false,
			Message: "Cannot cancel past appointments",
		}, nil
	}

	// Execute Transaction: DELETE schedule + DECREMENT used_sessions + INSERT log
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Delete TRAINING_SCHEDULE
	err = u.scheduleRepo.DeleteAppointment(ctx, tx, req.AppointmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete appointment: %w", err)
	}

	// 2. Decrement CUSTOMER_SESSION (used_sessions - 1)
	err = u.sessionRepo.DecrementUsedSessions(ctx, tx, appointment.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to decrement used sessions: %w", err)
	}

	// 3. Insert CUSTOMER_LOG (log_type = 'CANCEL_SESSION')
	err = u.customerLogRepo.CreateCustomerLog(ctx, tx, req.CustomerUsername, "CANCEL_SESSION")
	if err != nil {
		return nil, fmt.Errorf("failed to create customer log: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Get updated session info
	activeSession, err := u.sessionRepo.GetActiveSessionByCustomer(ctx, req.CustomerUsername)
	if err != nil {
		return nil, fmt.Errorf("failed to get active session: %w", err)
	}

	var remainingSessions int32
	if activeSession != nil {
		remainingSessions = activeSession.TotalSessions - activeSession.UsedSessions
	}

	// Return success response
	return &responses.CancelAppointmentResponse{
		Success:          true,
		Message:          "Appointment canceled successfully",
		AppointmentID:    req.AppointmentID,
		CustomerUsername: req.CustomerUsername,
		StartTime:        appointment.StartTime,
		EndTime:          appointment.EndTime,
		SessionID:        appointment.SessionID,
		RemainingSession: remainingSessions,
	}, nil
}
