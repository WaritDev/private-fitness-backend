package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/domain/requests"
	"github.com/WaritDev/private-fitness-backend/domain/responses"
)

// TrainerUseCase handles trainer-related business logic
type TrainerUseCase struct {
	trainerRepo repositories.TrainerRepository
}

// ProvideTrainerUseCase creates a new TrainerUseCase
func ProvideTrainerUseCase(
	trainerRepo repositories.TrainerRepository,
) *TrainerUseCase {
	return &TrainerUseCase{
		trainerRepo: trainerRepo,
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

	// Parse time strings to time.Time
	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return &responses.AddWorkingTimeResponse{
			Status:  "error",
			Message: "Invalid start time format. Expected HH:MM",
		}, nil
	}

	endTime, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		return &responses.AddWorkingTimeResponse{
			Status:  "error",
			Message: "Invalid end time format. Expected HH:MM",
		}, nil
	}

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
