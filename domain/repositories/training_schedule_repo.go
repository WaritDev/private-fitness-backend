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

	// Use Case 3P: Manage Day-Offs

	// GetTrainerDayOffs - Q3P.1: ดึงรายการวันหยุดทั้งหมดของ Trainer
	GetTrainerDayOffs(ctx context.Context, trainerUsername string) ([]DayOffInfo, error)

	// CheckDayOffDuplicate - Q3P.2: ตรวจสอบว่ามีวันหยุดซ้ำในวันนั้นหรือไม่
	CheckDayOffDuplicate(ctx context.Context, trainerUsername string, dayOffDate time.Time) (int64, error)

	// CheckDayOffAppointmentOverlap - Q3P.3: ตรวจสอบว่ามีนัดหมายใน day-off หรือไม่
	CheckDayOffAppointmentOverlap(ctx context.Context, trainerUsername string, startTime, endTime time.Time) (int64, error)

	// CreateDayOff - Q3P.4: สร้างวันหยุดใหม่
	CreateDayOff(ctx context.Context, trainerUsername string, startTime, endTime time.Time) error

	// DeleteDayOff - Q3P.5: ลบวันหยุด
	DeleteDayOff(ctx context.Context, scheduleID int32) error

	// Check-in Flow: Trainer Calendar
	// GetCustomerScheduleForToday - หา schedule ของลูกค้าในวันนี้ (สำหรับ check-in)
	GetCustomerScheduleForToday(ctx context.Context, customerUsername string) (*CustomerScheduleInfo, error)

	// GetTrainerAppointmentsWithPendingCheckIns - ดึง appointments พร้อม pending check-ins
	GetTrainerAppointmentsWithPendingCheckIns(ctx context.Context, trainerUsername string) ([]AppointmentWithCheckInInfo, error)
}

// CustomerScheduleInfo - ข้อมูล schedule ของลูกค้าในวันนี้
type CustomerScheduleInfo struct {
	ID               int32
	TrainerUsername  string
	CustomerUsername string
	SessionID        int32
	StartTime        time.Time
	EndTime          time.Time
}

// AppointmentWithCheckInInfo - ข้อมูล appointment พร้อม check-in status
type AppointmentWithCheckInInfo struct {
	ScheduleID        int32
	CustomerUsername  string
	CustomerFirstName string
	CustomerLastName  string
	StartTime         time.Time
	EndTime           time.Time
	SessionID         int32
	TotalSessions     int32
	UsedSessions      int32
	CheckinStatus     string // "PENDING", "CONFIRMED", "NONE"
	CheckinLogID      int32
	CheckinTime       time.Time
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
	ID               int32
	StartTime        time.Time
	EndTime          time.Time
	CustomerUsername string
	CheckinStatus    string // "PENDING", "CONFIRMED", "NONE"
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

// DayOffInfo - ข้อมูลวันหยุดของ Trainer (Use Case 3P)
type DayOffInfo struct {
	ScheduleID int32
	StartTime  time.Time
	EndTime    time.Time
}
