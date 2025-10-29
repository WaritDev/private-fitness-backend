-- name: ListUsers :many
SELECT id, username, email, created_at
FROM users
ORDER BY id;

-- name: CreateUser :execresult
INSERT INTO users (username, email)
VALUES (?, ?);