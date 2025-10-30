-- name: ListUsers :many
SELECT username, email, created_at
FROM users
ORDER BY username;

-- name: CreateUser :exec
INSERT INTO users (username, password, email, role, first_name, last_name, gender, date_of_birth, phone_number, gmail, specialty, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);


-- name: GetUserByUsername :one
SELECT username, password, role, first_name, last_name FROM users WHERE username = $1 LIMIT 1


