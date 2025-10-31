-- name: ListStaffs :many
SELECT
    username,
    role,
    first_name,
    last_name,
    gender,
    date_of_birth,
    phone_number,
    gmail,
    specialty,
    is_active
FROM users
WHERE role IN ('TRAINER','SALES','MANAGER','ADMIN')
ORDER BY is_active DESC, role ASC, first_name ASC, last_name ASC;

-- name: CheckPhoneExists :one
SELECT COUNT(phone_number) AS count
FROM users
WHERE phone_number = ?;

-- name: CheckGmailExists :one
SELECT COUNT(gmail) AS count
FROM users
WHERE LOWER(gmail) = LOWER(?);

-- name: CreateStaff :exec
INSERT INTO users (
    username, password, role, first_name, last_name,
    gender, date_of_birth, phone_number, gmail, specialty, is_active
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, TRUE);

-- name: GetStaffByUsername :one
SELECT
    username,
    role,
    first_name,
    last_name,
    gender,
    date_of_birth,
    phone_number,
    gmail,
    specialty,
    is_active
FROM users
WHERE username = ?;

-- name: UpdateStaffNoPassword :exec
UPDATE users
SET role         = ?,
    first_name   = ?,
    last_name    = ?,
    gender       = ?,
    date_of_birth= ?,
    phone_number = ?,
    gmail        = ?,
    specialty    = ?,
    is_active    = ?,
    updated_at   = CURRENT_TIMESTAMP
WHERE username = ?;

-- name: UpdateStaffWithPassword :exec
UPDATE users
SET password     = ?,
    role         = ?,
    first_name   = ?,
    last_name    = ?,
    gender       = ?,
    date_of_birth= ?,
    phone_number = ?,
    gmail        = ?,
    specialty    = ?,
    is_active    = ?,
    updated_at   = CURRENT_TIMESTAMP
WHERE username = ?;

-- name: CheckPhoneExistsExceptUsername :one
SELECT COUNT(phone_number) AS count
FROM users
WHERE phone_number = ?
    AND username != ?;

-- name: CheckGmailExistsExceptUsername :one
SELECT COUNT(username) as count
FROM users
WHERE LOWER(gmail) = LOWER(?)
    AND username != ?;

-- name: GetUserRole :one
SELECT role
FROM users
WHERE username = ?;

-- name: DeleteTrainerAvailabilityByTrainer :exec
DELETE FROM training_availabilities
WHERE trainer_username = ?;

-- name: DeleteTrainingScheduleByTrainer :exec
DELETE FROM training_schedules
WHERE trainer_username = ?;

-- name: DeleteCustomerSessionByTrainer :exec
DELETE FROM customer_sessions
WHERE trainer_username = ?;

-- name: DeleteCustomerDurationBySales :exec
DELETE FROM customer_durations
WHERE sales_username = ?;

-- name: DeleteCustomerSessionBySales :exec
DELETE FROM customer_sessions
WHERE sales_username = ?;

-- name: DeleteUserByUsername :exec
DELETE FROM users
WHERE username = ?;