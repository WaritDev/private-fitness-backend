-- name: CreateCustomerLog :exec
INSERT INTO customer_logs (
  customer_username,
  log_type
) VALUES (
  ?, ?
);

-- name: ListCustomerLogs :many
SELECT
  cl.id                    AS log_id,
  cl.customer_username     AS customer_username,
  u.first_name             AS customer_first_name,
  u.last_name              AS customer_last_name,
  cl.created_at            AS created_at,
  cl.log_type              AS log_type
FROM customer_logs cl
JOIN customers c ON c.username = cl.customer_username
JOIN users     u ON u.username = c.username
ORDER BY cl.created_at DESC, cl.id DESC;

-- name: GetCustomerLogByID :one
SELECT
  cl.id,
  cl.customer_username,
  u.first_name  AS customer_first_name,
  u.last_name   AS customer_last_name,
  cl.created_at,
  cl.log_type
FROM customer_logs cl
JOIN users u ON u.username = cl.customer_username
WHERE cl.id = ?
LIMIT 1;

-- name: UpdateCustomerLogByID :execresult
UPDATE customer_logs
SET created_at = ?,          -- DATETIME/TIMESTAMP
    log_type   = ?           -- ENUM('CHECK_IN','CHECK_OUT','BOOK_SESSION','CANCEL_SESSION')
WHERE id = ?;

-- name: DeleteCustomerLogByID :execresult
DELETE FROM customer_logs
WHERE id = ?;