-- name: ListUsers :many
SELECT id, username, email, created_at FROM users ORDER BY id;

-- name: CreateUser :one
INSERT INTO users (username, email)
VALUES ($1, $2)
RETURNING id, username, email, created_at;