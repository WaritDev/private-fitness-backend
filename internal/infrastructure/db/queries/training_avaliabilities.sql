-- name: GetTrainingAvaliabilitiesByTrainerUsername :many
SELECT 
  trainer_username,
  day_of_week,
  start_time,
  end_time
FROM training_availabilities
WHERE trainer_username = ?
ORDER BY day_of_week, start_time ASC;
