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
