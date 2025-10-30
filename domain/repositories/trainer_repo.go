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
}
