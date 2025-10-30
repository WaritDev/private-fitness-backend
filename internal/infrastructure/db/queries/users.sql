-- name: ListUsers :many
SELECT username, email, created_at
FROM users
ORDER BY username;

-- name: CreateUser :exec
INSERT INTO users (username, password, email, role, first_name, last_name, gender, date_of_birth, phone_number, gmail, specialty, is_active)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetUserByUsername :one
SELECT username, password, role, first_name, last_name
FROM users
WHERE username = ?
LIMIT 1;

-- name: CheckUsernameExists :one
SELECT COUNT(username) as count
FROM users
WHERE username = ?;

-- name: CheckEmailExists :one
SELECT COUNT(email) as count
FROM users
WHERE email = ?;

-- name: CheckPhoneNumberExists :one
SELECT COUNT(phone_number) as count
FROM users
WHERE phone_number = ?;

