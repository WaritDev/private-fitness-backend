package sql

import (
	"context"
	"database/sql"
	"time"

	"github.com/WaritDev/private-fitness-backend/domain/repositories"
	"github.com/WaritDev/private-fitness-backend/internal/infrastructure/db/dbmodel"
)

// TrainerRepository implements domain/repositories.TrainerRepository using sqlc
type TrainerRepository struct {
	q *dbmodel.Queries
}

// ProvideTrainerRepository creates a new TrainerRepository
func ProvideTrainerRepository(q *dbmodel.Queries) repositories.TrainerRepository {
	return &TrainerRepository{q: q}
}

// FindAvailableTrainers finds trainers available on specific day and time
func (r *TrainerRepository) FindAvailableTrainers(ctx context.Context, dayOfWeek string, startTime, endTime time.Time) ([]repositories.TrainerInfo, error) {
	rows, err := r.q.FindAvailableTrainers(ctx, dbmodel.FindAvailableTrainersParams{
		DayOfWeek: dbmodel.TrainingAvailabilitiesDayOfWeek(dayOfWeek),
		TIME:      startTime,
		TIME_2:    endTime,
	})
	if err != nil {
		return nil, err
	}

	result := make([]repositories.TrainerInfo, len(rows))
	for i, row := range rows {
		result[i] = repositories.TrainerInfo{
			Username:  row.Username,
			FirstName: row.FirstName,
			LastName:  row.LastName,
			CreatedAt: row.CreatedAt,
		}
	}
	return result, nil
}

// CountAppointmentsOnDate counts appointments for a trainer on specific date
func (r *TrainerRepository) CountAppointmentsOnDate(ctx context.Context, trainerUsername string, date time.Time) (int64, error) {
	count, err := r.q.CountAppointmentsOnDate(ctx, dbmodel.CountAppointmentsOnDateParams{
		TrainerUsername: sql.NullString{String: trainerUsername, Valid: true},
		DATE:            date,
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CheckScheduleOverlap checks if trainer has overlapping appointments
func (r *TrainerRepository) CheckScheduleOverlap(ctx context.Context, trainerUsername string, startTime, endTime time.Time) (bool, error) {
	count, err := r.q.CheckScheduleOverlap(ctx, dbmodel.CheckScheduleOverlapParams{
		TrainerUsername: sql.NullString{String: trainerUsername, Valid: true},
		StartTime:       endTime,   // start_time < ? (check if any appointment starts before our end)
		EndTime:         startTime, // end_time > ? (check if any appointment ends after our start)
	})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ListAllTrainers lists all active trainers
func (r *TrainerRepository) ListAllTrainers(ctx context.Context) ([]repositories.TrainerInfo, error) {
	rows, err := r.q.ListAllTrainers(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]repositories.TrainerInfo, len(rows))
	for i, row := range rows {
		result[i] = repositories.TrainerInfo{
			Username:  row.Username,
			FirstName: row.FirstName,
			LastName:  row.LastName,
		}
	}
	return result, nil
}

// Use Case 1P: Manage Working Hours

// GetTrainerAvailability gets all working hours for a trainer (Q1P.1)
func (r *TrainerRepository) GetTrainerAvailability(ctx context.Context, trainerUsername string) ([]repositories.TrainerAvailability, error) {
	rows, err := r.q.GetTrainerAvailability(ctx, trainerUsername)
	if err != nil {
		return nil, err
	}

	result := make([]repositories.TrainerAvailability, len(rows))
	for i, row := range rows {
		result[i] = repositories.TrainerAvailability{
			ID:              row.ID,
			TrainerUsername: row.TrainerUsername,
			DayOfWeek:       string(row.DayOfWeek),
			StartTime:       row.StartTime,
			EndTime:         row.EndTime,
		}
	}
	return result, nil
}

// CheckTimeOverlap checks if new time slot overlaps with existing ones (Q1P.2)
func (r *TrainerRepository) CheckTimeOverlap(ctx context.Context, trainerUsername string, dayOfWeek string, startTime, endTime time.Time) (int64, error) {
	count, err := r.q.CheckTimeOverlap(ctx, dbmodel.CheckTimeOverlapParams{
		TrainerUsername: trainerUsername,
		DayOfWeek:       dbmodel.TrainingAvailabilitiesDayOfWeek(dayOfWeek),
		TIME:            startTime,
		TIME_2:          endTime,
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CreateTrainerAvailability creates new working time slot (Q1P.3)
func (r *TrainerRepository) CreateTrainerAvailability(ctx context.Context, trainerUsername string, dayOfWeek string, startTime, endTime time.Time) error {
	err := r.q.CreateTrainerAvailability(ctx, dbmodel.CreateTrainerAvailabilityParams{
		TrainerUsername: trainerUsername,
		DayOfWeek:       dbmodel.TrainingAvailabilitiesDayOfWeek(dayOfWeek),
		StartTime:       startTime,
		EndTime:         endTime,
	})
	return err
}

// UpdateTrainerAvailability updates existing working time slot (Q1P.4)
func (r *TrainerRepository) UpdateTrainerAvailability(ctx context.Context, id int32, dayOfWeek string, startTime, endTime time.Time) error {
	err := r.q.UpdateTrainerAvailability(ctx, dbmodel.UpdateTrainerAvailabilityParams{
		DayOfWeek: dbmodel.TrainingAvailabilitiesDayOfWeek(dayOfWeek),
		StartTime: startTime,
		EndTime:   endTime,
		ID:        id,
	})
	return err
}

// DeleteTrainerAvailability deletes a working time slot (Q1P.5)
func (r *TrainerRepository) DeleteTrainerAvailability(ctx context.Context, id int32) error {
	err := r.q.DeleteTrainerAvailability(ctx, id)
	return err
}

// GetTrainerAvailabilityByID gets a single working hour by ID (for validation)
func (r *TrainerRepository) GetTrainerAvailabilityByID(ctx context.Context, id int32) (*repositories.TrainerAvailability, error) {
	// We'll need to query by ID - use GetTrainerAvailability and filter
	// This is a workaround since we don't have a specific query for single ID
	// In production, you might want to add a dedicated SQL query for this

	// For now, we'll return nil to indicate we need to implement the query
	// The use case layer will handle this by checking ownership differently
	return nil, sql.ErrNoRows
}
