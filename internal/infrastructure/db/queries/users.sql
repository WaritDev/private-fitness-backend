-- name: ListUsers :many
SELECT username, gmail, created_at
FROM users
ORDER BY username;

-- name: CreateUser :exec
INSERT INTO users (username, password, role, first_name, last_name, gender, date_of_birth, phone_number, gmail, specialty, is_active)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetUserByUsername :one
SELECT username, password, role, first_name, last_name
FROM users
WHERE username = ?
LIMIT 1;

-- name: CheckUsernameExists :one
SELECT COUNT(username) as count
FROM users
WHERE username = ?;

-- name: CheckGmailExistsUser :one
SELECT COUNT(gmail) as count
FROM users
WHERE gmail = ?;

-- name: CheckPhoneNumberExistsUser :one
SELECT COUNT(phone_number) as count
FROM users
WHERE phone_number = ?;

-- name: UpdateUserLoginTimestamp :exec
-- Q0S.2: Update updated_at when user logs in (track last active time)
UPDATE users
SET updated_at = CURRENT_TIMESTAMP
WHERE username = ?;
