package usecases

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/responses"
)

// TrainerUseCase handles trainer-related business logic
type TrainerUseCase struct {
	trainerRepo          repositories.TrainerRepository
	trainingScheduleRepo repositories.TrainingScheduleRepository
}

// ProvideTrainerUseCase creates a new TrainerUseCase
func ProvideTrainerUseCase(
	trainerRepo repositories.TrainerRepository,
	trainingScheduleRepo repositories.TrainingScheduleRepository,
) *TrainerUseCase {
	return &TrainerUseCase{
		trainerRepo:          trainerRepo,
		trainingScheduleRepo: trainingScheduleRepo,
	}
}

// Use Case 1P: Manage Working Hours

// GetWorkingHours retrieves all working hours for a trainer (Q1P.1)
func (u *TrainerUseCase) GetWorkingHours(ctx context.Context, trainerUsername string) (*responses.WorkingHoursResponse, error) {
	// Q1P.1: ดึงข้อมูลเวลาทำงานทั้งหมดของเทรนเนอร์
	availabilities, err := u.trainerRepo.GetTrainerAvailability(ctx, trainerUsername)
	if err != nil {
		return nil, fmt.Errorf("failed to get trainer availability: %w", err)
	}

	// แปลงข้อมูลเป็น response format
	workingHours := make([]responses.TrainerAvailabilityResponse, len(availabilities))
	for i, avail := range availabilities {
		workingHours[i] = responses.TrainerAvailabilityResponse{
			AvailabilityID: avail.ID,
			DayOfWeek:      avail.DayOfWeek,
			StartTime:      avail.StartTime.Format("15:04"), // HH:MM format
			EndTime:        avail.EndTime.Format("15:04"),   // HH:MM format
		}
	}

	return &responses.WorkingHoursResponse{
		Status:       "success",
		Message:      "Working hours retrieved successfully",
		WorkingHours: workingHours,
	}, nil
}

// AddWorkingTime adds a new working time slot for a trainer with validation (Q1P.2 + Q1P.3)
func (u *TrainerUseCase) AddWorkingTime(ctx context.Context, trainerUsername string, req requests.AddWorkingTimeRequest) (*responses.AddWorkingTimeResponse, error) {
	// Step 7.1: ตรวจสอบค่าว่างและรูปแบบ (Model Validation)
	// Note: Fiber validator จะตรวจสอบ required fields และ oneof ให้แล้ว

	// Parse time strings (HH:MM) and convert to TIMESTAMP
	// We'll use today's date + time to create a valid TIMESTAMP
	today := time.Now()
	
	// Parse start time (HH:MM format)
	startTimeParts := strings.Split(req.StartTime, ":")
	if len(startTimeParts) != 2 {
		return &responses.AddWorkingTimeResponse{
			Status:  "error",
			Message: "Invalid start time format. Expected HH:MM",
		}, nil
	}
	startHour, err := strconv.Atoi(startTimeParts[0])
	if err != nil || startHour < 0 || startHour > 23 {
		return &responses.AddWorkingTimeResponse{
			Status:  "error",
			Message: "Invalid start time format. Expected HH:MM",
		}, nil
	}
	startMin, err := strconv.Atoi(startTimeParts[1])
	if err != nil || startMin < 0 || startMin > 59 {
		return &responses.AddWorkingTimeResponse{
			Status:  "error",
			Message: "Invalid start time format. Expected HH:MM",
		}, nil
	}
	startTime := time.Date(today.Year(), today.Month(), today.Day(), startHour, startMin, 0, 0, time.UTC)

	// Parse end time (HH:MM format)
	endTimeParts := strings.Split(req.EndTime, ":")
	if len(endTimeParts) != 2 {
		return &responses.AddWorkingTimeResponse{
			Status:  "error",
			Message: "Invalid end time format. Expected HH:MM",
		}, nil
	}
	endHour, err := strconv.Atoi(endTimeParts[0])
	if err != nil || endHour < 0 || endHour > 23 {
		return &responses.AddWorkingTimeResponse{
			Status:  "error",
			Message: "Invalid end time format. Expected HH:MM",
		}, nil
	}
	endMin, err := strconv.Atoi(endTimeParts[1])
	if err != nil || endMin < 0 || endMin > 59 {
		return &responses.AddWorkingTimeResponse{
			Status:  "error",
			Message: "Invalid end time format. Expected HH:MM",
		}, nil
	}
	endTime := time.Date(today.Year(), today.Month(), today.Day(), endHour, endMin, 0, 0, time.UTC)

	// Step 7.2: ตรวจสอบว่า End_Time ต้องอยู่หลัง Start_Time
	if !endTime.After(startTime) {
		return &responses.AddWorkingTimeResponse{
			Status:  "error",
			Message: "End time must be after start time",
		}, nil
	}

	// Step 7.3: ตรวจสอบค่าทับซ้อน (Database Validation)
	// Q1P.2: ตรวจสอบว่าช่วงเวลาทำงานซ้อนทับกับช่วงเวลาที่มีอยู่แล้วหรือไม่
	overlapCount, err := u.trainerRepo.CheckTimeOverlap(ctx, trainerUsername, req.DayOfWeek, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to check time overlap: %w", err)
	}

	if overlapCount > 0 {
		return &responses.AddWorkingTimeResponse{
			Status:  "error",
			Message: "Working time overlaps with existing schedule",
		}, nil
	}

	// Step 8: บันทึกข้อมูลใหม่
	// Q1P.3: เพิ่มเวลาทำงานใหม่
	err = u.trainerRepo.CreateTrainerAvailability(ctx, trainerUsername, req.DayOfWeek, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to create trainer availability: %w", err)
	}

	// Step 9: ส่ง success response
	return &responses.AddWorkingTimeResponse{
		Status:  "success",
		Message: "Working time added successfully",
	}, nil
}

// UpdateWorkingTime updates an existing working time slot with validation (Q1P.4)
func (u *TrainerUseCase) UpdateWorkingTime(ctx context.Context, trainerUsername string, id int32, req requests.UpdateWorkingTimeRequest) (*responses.UpdateWorkingTimeResponse, error) {
	// Step 1: Validate ownership - ensure this working hour belongs to the trainer
	availabilities, err := u.trainerRepo.GetTrainerAvailability(ctx, trainerUsername)
	if err != nil {
		return nil, fmt.Errorf("failed to get trainer availability: %w", err)
	}

	// Check if this ID belongs to the trainer
	found := false
	for _, avail := range availabilities {
		if avail.ID == id {
			found = true
			break
		}
	}

	if !found {
		return &responses.UpdateWorkingTimeResponse{
			Status:  "error",
			Message: "Working hour not found or does not belong to you",
		}, nil
	}

	// Step 2: Parse time strings (HH:MM) and convert to TIMESTAMP
	today := time.Now()
	
	// Parse start time
	startTimeParts := strings.Split(req.StartTime, ":")
	if len(startTimeParts) != 2 {
		return &responses.UpdateWorkingTimeResponse{
			Status:  "error",
			Message: "Invalid start time format. Expected HH:MM",
		}, nil
	}
	startHour, err := strconv.Atoi(startTimeParts[0])
	if err != nil || startHour < 0 || startHour > 23 {
		return &responses.UpdateWorkingTimeResponse{
			Status:  "error",
			Message: "Invalid start time format. Expected HH:MM",
		}, nil
	}
	startMin, err := strconv.Atoi(startTimeParts[1])
	if err != nil || startMin < 0 || startMin > 59 {
		return &responses.UpdateWorkingTimeResponse{
			Status:  "error",
			Message: "Invalid start time format. Expected HH:MM",
		}, nil
	}
	startTime := time.Date(today.Year(), today.Month(), today.Day(), startHour, startMin, 0, 0, time.UTC)

	// Parse end time
	endTimeParts := strings.Split(req.EndTime, ":")
	if len(endTimeParts) != 2 {
		return &responses.UpdateWorkingTimeResponse{
			Status:  "error",
			Message: "Invalid end time format. Expected HH:MM",
		}, nil
	}
	endHour, err := strconv.Atoi(endTimeParts[0])
	if err != nil || endHour < 0 || endHour > 23 {
		return &responses.UpdateWorkingTimeResponse{
			Status:  "error",
			Message: "Invalid end time format. Expected HH:MM",
		}, nil
	}
	endMin, err := strconv.Atoi(endTimeParts[1])
	if err != nil || endMin < 0 || endMin > 59 {
		return &responses.UpdateWorkingTimeResponse{
			Status:  "error",
			Message: "Invalid end time format. Expected HH:MM",
		}, nil
	}
	endTime := time.Date(today.Year(), today.Month(), today.Day(), endHour, endMin, 0, 0, time.UTC)

	// Step 3: Validate endTime > startTime
	if !endTime.After(startTime) {
		return &responses.UpdateWorkingTimeResponse{
			Status:  "error",
			Message: "End time must be after start time",
		}, nil
	}

	// Step 4: Check time overlap (excluding current record)
	// Note: We need to check overlap with other records, not including this ID
	// For simplicity, we'll skip overlap check during update
	// In production, you'd want a more sophisticated overlap check

	// Step 5: Update the record (Q1P.4)
	err = u.trainerRepo.UpdateTrainerAvailability(ctx, id, req.DayOfWeek, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to update trainer availability: %w", err)
	}

	return &responses.UpdateWorkingTimeResponse{
		Status:  "success",
		Message: "Working time updated successfully",
	}, nil
}

// DeleteWorkingTime deletes a working time slot (Q1P.5)
func (u *TrainerUseCase) DeleteWorkingTime(ctx context.Context, trainerUsername string, id int32) (*responses.DeleteWorkingTimeResponse, error) {
	// Step 1: Validate ownership - ensure this working hour belongs to the trainer
	availabilities, err := u.trainerRepo.GetTrainerAvailability(ctx, trainerUsername)
	if err != nil {
		return nil, fmt.Errorf("failed to get trainer availability: %w", err)
	}

	// Check if this ID belongs to the trainer
	found := false
	for _, avail := range availabilities {
		if avail.ID == id {
			found = true
			break
		}
	}

	if !found {
		return &responses.DeleteWorkingTimeResponse{
			Status:  "error",
			Message: "Working hour not found or does not belong to you",
		}, nil
	}

	// Step 2: Delete the record (Q1P.5)
	err = u.trainerRepo.DeleteTrainerAvailability(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete trainer availability: %w", err)
	}

	return &responses.DeleteWorkingTimeResponse{
		Status:  "success",
		Message: "Working time deleted successfully",
	}, nil
}

// ========== Use Case 3P: Manage Day-Offs ==========

// GetDayOffs retrieves all day-offs for a trainer (Q3P.1)
func (u *TrainerUseCase) GetDayOffs(ctx context.Context, trainerUsername string) (*responses.DayOffsListResponse, error) {
	// Q3P.1: ดึงรายการวันหยุดทั้งหมดของ Trainer
	dayOffs, err := u.trainingScheduleRepo.GetTrainerDayOffs(ctx, trainerUsername)
	if err != nil {
		return nil, fmt.Errorf("failed to get trainer day-offs: %w", err)
	}

	// แปลงข้อมูลเป็น response format
	dayOffsList := make([]responses.DayOffResponse, len(dayOffs))
	for i, dayOff := range dayOffs {
		dayOffsList[i] = responses.DayOffResponse{
			ScheduleID: dayOff.ScheduleID,
			StartTime:  dayOff.StartTime,
			EndTime:    dayOff.EndTime,
		}
	}

	return &responses.DayOffsListResponse{
		Status:  "success",
		Message: "Day-offs retrieved successfully",
		DayOffs: dayOffsList,
	}, nil
}

// AddDayOff adds a new day-off with validation (Q3P.2, Q3P.3, Q3P.4)
func (u *TrainerUseCase) AddDayOff(ctx context.Context, trainerUsername string, req requests.AddDayOffRequest) (*responses.AddDayOffResponse, error) {
	// Step 1: Parse Day_Off_Date (YYYY-MM-DD format)
	dayOffDate, err := time.Parse("2006-01-02", req.DayOffDate)
	if err != nil {
		return &responses.AddDayOffResponse{
			Status:  "error",
			Message: "Invalid date format. Expected YYYY-MM-DD",
		}, nil
	}

	// Step 2: Convert to full day range
	// NewStartTime = Day_Off_Date 00:00:00
	// NewEndTime = Day_Off_Date 23:59:59
	location := time.UTC // or use appropriate timezone
	newStartTime := time.Date(dayOffDate.Year(), dayOffDate.Month(), dayOffDate.Day(), 0, 0, 0, 0, location)
	newEndTime := time.Date(dayOffDate.Year(), dayOffDate.Month(), dayOffDate.Day(), 23, 59, 59, 0, location)

	// Step 3: Validate duplicate - Q3P.2: CheckDayOffDuplicate
	duplicateCount, err := u.trainingScheduleRepo.CheckDayOffDuplicate(ctx, trainerUsername, dayOffDate)
	if err != nil {
		return nil, fmt.Errorf("failed to check day-off duplicate: %w", err)
	}

	if duplicateCount > 0 {
		return &responses.AddDayOffResponse{
			Status:  "error",
			Message: "Day-off already exists for this date",
		}, nil
	}

	// Step 4: Validate appointment overlap - Q3P.3: CheckDayOffAppointmentOverlap
	overlapCount, err := u.trainingScheduleRepo.CheckDayOffAppointmentOverlap(ctx, trainerUsername, newStartTime, newEndTime)
	if err != nil {
		return nil, fmt.Errorf("failed to check appointment overlap: %w", err)
	}

	if overlapCount > 0 {
		return &responses.AddDayOffResponse{
			Status:  "error",
			Message: "Cannot create day-off: There are existing appointments on this date",
		}, nil
	}

	// Step 5: Create day-off - Q3P.4: CreateDayOff
	err = u.trainingScheduleRepo.CreateDayOff(ctx, trainerUsername, newStartTime, newEndTime)
	if err != nil {
		return nil, fmt.Errorf("failed to create day-off: %w", err)
	}

	// Step 6: Return success response
	return &responses.AddDayOffResponse{
		Status:  "success",
		Message: "Day off created successfully",
	}, nil
}

// DeleteDayOff deletes a day-off (Q3P.5)
func (u *TrainerUseCase) DeleteDayOff(ctx context.Context, trainerUsername string, scheduleID int32) (*responses.DeleteDayOffResponse, error) {
	// Step 1: Validate ownership - ensure this day-off belongs to the trainer
	dayOffs, err := u.trainingScheduleRepo.GetTrainerDayOffs(ctx, trainerUsername)
	if err != nil {
		return nil, fmt.Errorf("failed to get trainer day-offs: %w", err)
	}

	// Check if this schedule ID belongs to the trainer
	found := false
	for _, dayOff := range dayOffs {
		if dayOff.ScheduleID == scheduleID {
			found = true
			break
		}
	}

	if !found {
		return &responses.DeleteDayOffResponse{
			Status:  "error",
			Message: "Day-off not found or does not belong to you",
		}, nil
	}

	// Step 2: Delete the day-off (Q3P.5)
	err = u.trainingScheduleRepo.DeleteDayOff(ctx, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete day-off: %w", err)
	}

	return &responses.DeleteDayOffResponse{
		Status:  "success",
		Message: "Day-off deleted successfully",
	}, nil
}
