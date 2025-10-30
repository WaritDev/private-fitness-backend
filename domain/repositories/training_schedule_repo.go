package repositories

import (
	"context"
	"time"
)

// TrainingScheduleRepository - Repository for training schedules and availability
type TrainingScheduleRepository interface {
	// GetTrainerAvailability - ดึงเวลาทำงานประจำสัปดาห์ของ Trainer
	GetTrainerAvailability(ctx context.Context, trainerUsername string) ([]TrainerAvailabilityInfo, error)

	// GetDayOffSchedules - ดึงวันหยุดหรือช่วงเวลาที่ไม่รับนัด (DAY_OFF)
	GetDayOffSchedules(ctx context.Context, trainerUsername string, calendarStart, calendarEnd time.Time) ([]ScheduleTimeSlot, error)

	// GetAppointmentSchedules - ดึงนัดที่ถูกจองแล้ว (APPOINTMENT)
	GetAppointmentSchedules(ctx context.Context, trainerUsername string, calendarStart, calendarEnd time.Time) ([]AppointmentInfo, error)

	// CheckTimeSlotAvailability - Q3C.5: ตรวจสอบว่าช่วงเวลายังว่างอยู่
	CheckTimeSlotAvailability(ctx context.Context, trainerUsername string, startTime, endTime time.Time) (int64, error)

	// BookAppointment - Q3C.6: บันทึกการจอง (Transaction: INSERT schedule + UPDATE used_sessions + INSERT log)
	BookAppointment(ctx context.Context, params BookAppointmentParams) error

	// GetAppointmentById - ดึงข้อมูลการจองตาม ID
	GetAppointmentById(ctx context.Context, appointmentID int32) (*AppointmentDetail, error)

	// DeleteAppointment - ลบการจอง (ใช้ใน Transaction)
	DeleteAppointment(ctx context.Context, tx interface{}, appointmentID int32) error
}

// TrainerAvailabilityInfo - ข้อมูลเวลาทำงานประจำสัปดาห์
type TrainerAvailabilityInfo struct {
	TrainerUsername string
	DayOfWeek       string
	StartTime       time.Time
	EndTime         time.Time
}

// ScheduleTimeSlot - ช่วงเวลาของ schedule (ใช้สำหรับ DAY_OFF)
type ScheduleTimeSlot struct {
	StartTime time.Time
	EndTime   time.Time
}

// AppointmentInfo - ข้อมูลนัดหมาย
type AppointmentInfo struct {
	StartTime        time.Time
	EndTime          time.Time
	CustomerUsername string
}

// BookAppointmentParams - พารามิเตอร์สำหรับการจองนัด
type BookAppointmentParams struct {
	TrainerUsername  string
	CustomerUsername string
	SessionID        int32
	StartTime        time.Time
	EndTime          time.Time
}

// AppointmentDetail - รายละเอียดการจอง
type AppointmentDetail struct {
	ID               int32
	TrainerUsername  string
	CustomerUsername string
	SessionID        int32
	StartTime        time.Time
	EndTime          time.Time
	ScheduleType     string
}
