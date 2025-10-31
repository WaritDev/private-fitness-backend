package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

// TrainingScheduleRepository implements domain/repositories.TrainingScheduleRepository
type TrainingScheduleRepository struct {
	q *dbmodel.Queries
}

// ProvideTrainingScheduleRepository creates a new TrainingScheduleRepository
func ProvideTrainingScheduleRepository(q *dbmodel.Queries) repositories.TrainingScheduleRepository {
	return &TrainingScheduleRepository{q: q}
}

// GetTrainerAvailability - ดึงเวลาทำงานประจำสัปดาห์ของ Trainer
func (r *TrainingScheduleRepository) GetTrainerAvailability(ctx context.Context, trainerUsername string) ([]repositories.TrainerAvailabilityInfo, error) {
	rows, err := r.q.GetTrainingAvaliabilitiesByTrainerUsername(ctx, trainerUsername)
	if err != nil {
		return nil, err
	}

	result := make([]repositories.TrainerAvailabilityInfo, len(rows))
	for i, row := range rows {
		result[i] = repositories.TrainerAvailabilityInfo{
			TrainerUsername: row.TrainerUsername,
			DayOfWeek:       string(row.DayOfWeek),
			StartTime:       row.StartTime,
			EndTime:         row.EndTime,
		}
	}

	return result, nil
}

// GetDayOffSchedules - ดึงวันหยุดหรือช่วงเวลาที่ไม่รับนัด (DAY_OFF)
func (r *TrainingScheduleRepository) GetDayOffSchedules(ctx context.Context, trainerUsername string, calendarStart, calendarEnd time.Time) ([]repositories.ScheduleTimeSlot, error) {
	rows, err := r.q.GetDayOffSchedules(ctx, dbmodel.GetDayOffSchedulesParams{
		TrainerUsername: sql.NullString{String: trainerUsername, Valid: true},
		StartTime:       calendarEnd,
		EndTime:         calendarStart,
	})
	if err != nil {
		return nil, err
	}

	result := make([]repositories.ScheduleTimeSlot, len(rows))
	for i, row := range rows {
		result[i] = repositories.ScheduleTimeSlot{
			StartTime: row.StartTime,
			EndTime:   row.EndTime,
		}
	}

	return result, nil
}

// GetAppointmentSchedules - ดึงนัดที่ถูกจองแล้ว (APPOINTMENT)
func (r *TrainingScheduleRepository) GetAppointmentSchedules(ctx context.Context, trainerUsername string, calendarStart, calendarEnd time.Time) ([]repositories.AppointmentInfo, error) {
	rows, err := r.q.GetAppointmentSchedules(ctx, dbmodel.GetAppointmentSchedulesParams{
		TrainerUsername: sql.NullString{String: trainerUsername, Valid: true},
		StartTime:       calendarEnd,
		EndTime:         calendarStart,
	})
	if err != nil {
		return nil, err
	}

	result := make([]repositories.AppointmentInfo, len(rows))
	for i, row := range rows {
		var customerUser string
		if row.CustomerUsername.Valid {
			customerUser = row.CustomerUsername.String
		}

		result[i] = repositories.AppointmentInfo{
			ID:               row.ID, // เพิ่ม ID สำหรับการยกเลิกนัดและแสดงใน response
			StartTime:        row.StartTime,
			EndTime:          row.EndTime,
			CustomerUsername: customerUser,
		}
	}

	return result, nil
}

// CheckTimeSlotAvailability - Q3C.5: ตรวจสอบว่าช่วงเวลายังว่างอยู่
func (r *TrainingScheduleRepository) CheckTimeSlotAvailability(ctx context.Context, trainerUsername string, startTime, endTime time.Time) (int64, error) {
	count, err := r.q.CheckTimeSlotAvailability(ctx, dbmodel.CheckTimeSlotAvailabilityParams{
		TrainerUsername: sql.NullString{String: trainerUsername, Valid: true},
		StartTime:       endTime, // Parameters order: trainer, endTime, startTime
		EndTime:         startTime,
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// BookAppointment - Q3C.6: บันทึกการจอง (Transaction: INSERT schedule + UPDATE used_sessions + INSERT log)
func (r *TrainingScheduleRepository) BookAppointment(ctx context.Context, params repositories.BookAppointmentParams) error {
	// Note: Transaction จะถูกจัดการใน use case layer
	// Repository นี้จะถูกเรียกผ่าน transaction context
	return r.q.BookAppointment(ctx, dbmodel.BookAppointmentParams{
		TrainerUsername:  sql.NullString{String: params.TrainerUsername, Valid: true},
		CustomerUsername: sql.NullString{String: params.CustomerUsername, Valid: true},
		SessionID:        sql.NullInt32{Int32: params.SessionID, Valid: true},
		StartTime:        params.StartTime,
		EndTime:          params.EndTime,
	})
}

// GetAppointmentById - ดึงข้อมูลการจองตาม ID
func (r *TrainingScheduleRepository) GetAppointmentById(ctx context.Context, appointmentID int32) (*repositories.AppointmentDetail, error) {
	row, err := r.q.GetAppointmentById(ctx, appointmentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Appointment not found
		}
		return nil, err
	}

	var trainerUsername, customerUsername string
	if row.TrainerUsername.Valid {
		trainerUsername = row.TrainerUsername.String
	}
	if row.CustomerUsername.Valid {
		customerUsername = row.CustomerUsername.String
	}

	var sessionID int32
	if row.SessionID.Valid {
		sessionID = row.SessionID.Int32
	}

	return &repositories.AppointmentDetail{
		ID:               row.ID,
		TrainerUsername:  trainerUsername,
		CustomerUsername: customerUsername,
		SessionID:        sessionID,
		StartTime:        row.StartTime,
		EndTime:          row.EndTime,
		ScheduleType:     string(row.ScheduleType),
	}, nil
}

// DeleteAppointment - ลบการจอง (ใช้ใน Transaction)
func (r *TrainingScheduleRepository) DeleteAppointment(ctx context.Context, tx interface{}, appointmentID int32) error {
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	qtx := r.q.WithTx(sqlTx)
	return qtx.DeleteAppointment(ctx, appointmentID)
}

// ========== Use Case 3P: Manage Day-Offs ==========

// GetTrainerDayOffs - Q3P.1: ดึงรายการวันหยุดทั้งหมดของ Trainer
func (r *TrainingScheduleRepository) GetTrainerDayOffs(ctx context.Context, trainerUsername string) ([]repositories.DayOffInfo, error) {
	rows, err := r.q.GetTrainerDayOffs(ctx, sql.NullString{String: trainerUsername, Valid: true})
	if err != nil {
		return nil, err
	}

	result := make([]repositories.DayOffInfo, len(rows))
	for i, row := range rows {
		result[i] = repositories.DayOffInfo{
			ScheduleID: row.ID,
			StartTime:  row.StartTime,
			EndTime:    row.EndTime,
		}
	}

	return result, nil
}

// CheckDayOffDuplicate - Q3P.2: ตรวจสอบว่ามีวันหยุดซ้ำในวันนั้นหรือไม่
func (r *TrainingScheduleRepository) CheckDayOffDuplicate(ctx context.Context, trainerUsername string, dayOffDate time.Time) (int64, error) {
	count, err := r.q.CheckDayOffDuplicate(ctx, dbmodel.CheckDayOffDuplicateParams{
		TrainerUsername: sql.NullString{String: trainerUsername, Valid: true},
		StartTime:       dayOffDate,
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CheckDayOffAppointmentOverlap - Q3P.3: ตรวจสอบว่ามีนัดหมายใน day-off หรือไม่
func (r *TrainingScheduleRepository) CheckDayOffAppointmentOverlap(ctx context.Context, trainerUsername string, startTime, endTime time.Time) (int64, error) {
	count, err := r.q.CheckDayOffAppointmentOverlap(ctx, dbmodel.CheckDayOffAppointmentOverlapParams{
		TrainerUsername: sql.NullString{String: trainerUsername, Valid: true},
		EndTime:         endTime,   // WHERE start_time < ?
		StartTime:       startTime, // AND end_time > ?
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CreateDayOff - Q3P.4: สร้างวันหยุดใหม่
func (r *TrainingScheduleRepository) CreateDayOff(ctx context.Context, trainerUsername string, startTime, endTime time.Time) error {
	return r.q.CreateDayOff(ctx, dbmodel.CreateDayOffParams{
		TrainerUsername: sql.NullString{String: trainerUsername, Valid: true},
		StartTime:       startTime,
		EndTime:         endTime,
	})
}

// DeleteDayOff - Q3P.5: ลบวันหยุด
func (r *TrainingScheduleRepository) DeleteDayOff(ctx context.Context, scheduleID int32) error {
	return r.q.DeleteDayOff(ctx, scheduleID)
}
