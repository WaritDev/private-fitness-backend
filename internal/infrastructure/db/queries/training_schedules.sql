-- name: CreateTrainingSchedule :exec
INSERT INTO training_schedules (
  trainer_username,
  customer_username,
  session_id,
  start_time,
  end_time,
  schedule_type
) VALUES (
  ?, ?, ?, ?, ?, ?
);