-- name: CreateCustomerLog :exec
INSERT INTO customer_logs (
  customer_username,
  log_type
) VALUES (
  ?, ?
);