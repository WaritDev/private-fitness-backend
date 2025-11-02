package repositories

import (
	"context"
	"time"
)

// TrainerInfo - Trainer information from database
type TrainerInfo struct {
	Username  string
	FirstName string
	LastName  string
	CreatedAt time.Time
}

// TrainerAvailability - Trainer availability information
type TrainerAvailability struct {
	ID              int32
	TrainerUsername string
	DayOfWeek       string
	StartTime       time.Time
	EndTime         time.Time
}

// TrainerRepository - Repository interface for trainer operations
type TrainerRepository interface {
	// FindAvailableTrainers finds trainers available on specific day and time
	FindAvailableTrainers(ctx context.Context, dayOfWeek string, startTime, endTime time.Time) ([]TrainerInfo, error)

	// CountAppointmentsOnDate counts appointments for a trainer on specific date
	CountAppointmentsOnDate(ctx context.Context, trainerUsername string, date time.Time) (int64, error)

	// CheckScheduleOverlap checks if trainer has overlapping appointments
	CheckScheduleOverlap(ctx context.Context, trainerUsername string, startTime, endTime time.Time) (bool, error)

	// ListAllTrainers lists all active trainers
	ListAllTrainers(ctx context.Context) ([]TrainerInfo, error)

	// Use Case 1P: Manage Working Hours

	// GetTrainerAvailability gets all working hours for a trainer (Q1P.1)
	GetTrainerAvailability(ctx context.Context, trainerUsername string) ([]TrainerAvailability, error)

	// CheckTimeOverlap checks if new time slot overlaps with existing ones (Q1P.2)
	CheckTimeOverlap(ctx context.Context, trainerUsername string, dayOfWeek string, startTime, endTime time.Time) (int64, error)

	// CreateTrainerAvailability creates new working time slot (Q1P.3)
	CreateTrainerAvailability(ctx context.Context, trainerUsername string, dayOfWeek string, startTime, endTime time.Time) error

	// UpdateTrainerAvailability updates existing working time slot (Q1P.4)
	UpdateTrainerAvailability(ctx context.Context, id int32, dayOfWeek string, startTime, endTime time.Time) error

	// DeleteTrainerAvailability deletes a working time slot (Q1P.5)
	DeleteTrainerAvailability(ctx context.Context, id int32) error

	// GetTrainerAvailabilityByID gets a single working hour by ID (for validation)
	GetTrainerAvailabilityByID(ctx context.Context, id int32) (*TrainerAvailability, error)

	// CheckAvailabilityDayOffOverlap checks if there are working hours on a specific day of week
	CheckAvailabilityDayOffOverlap(ctx context.Context, trainerUsername string, dayOfWeek string) (int64, error)
}
