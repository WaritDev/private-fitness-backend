-- name: CreateCustomerDuration :exec
INSERT INTO customer_durations (
  customer_username,
  sales_username,
  product_id,
  purchase_date,
  start_date,
  end_date,
  price_paid,
  discount_amount,
  status
) VALUES (
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?
);

-- name: GetCustomerDurationsByUsername :many
SELECT
  id,
  customer_username,
  sales_username,
  product_id,
  purchase_date,
  start_date,
  end_date,
  price_paid,
  discount_amount,
  status,
  created_at,
  updated_at
FROM customer_durations
WHERE customer_username = ?
ORDER BY created_at DESC;

-- name: GetCustomerDurationById :one
SELECT
  id,
  customer_username,
  sales_username,
  product_id,
  purchase_date,
  start_date,
  end_date,
  price_paid,
  discount_amount,
  status,
  created_at,
  updated_at
FROM customer_durations
WHERE id = ?
LIMIT 1;
