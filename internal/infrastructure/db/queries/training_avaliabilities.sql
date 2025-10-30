-- name: GetTrainingAvaliabilitiesByTrainerUsername :many
SELECT 
  trainer_username,
  day_of_week,
  start_time,
  end_time
FROM training_availabilities
WHERE trainer_username = ?
ORDER BY day_of_week, start_time ASC;

-- Use Case 1P: Manage Working Hours

-- Q1P.1: Get Trainer Availability with ID (Get Working Hours)
-- name: GetTrainerAvailability :many
SELECT
  id,
  trainer_username,
  day_of_week,
  start_time,
  end_time
FROM training_availabilities
WHERE trainer_username = ?
ORDER BY FIELD(day_of_week, 'MONDAY', 'TUESDAY', 'WEDNESDAY', 'THURSDAY', 'FRIDAY', 'SATURDAY', 'SUNDAY'),
         start_time ASC;

-- Q1P.2: Check Time Overlap (Validation before adding)
-- name: CheckTimeOverlap :one
SELECT COUNT(id) AS overlapped_count
FROM training_availabilities
WHERE trainer_username = ?
  AND day_of_week = ?
  AND (
    (? < end_time AND ? > start_time)
  );

-- Q1P.3: Create Trainer Availability (Add Working Time)
-- name: CreateTrainerAvailability :exec
INSERT INTO training_availabilities (
  trainer_username,
  day_of_week,
  start_time,
  end_time
) VALUES (?, ?, ?, ?);
