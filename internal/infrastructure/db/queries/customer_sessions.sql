-- name: CreateCustomerSession :exec
INSERT INTO customer_sessions (
  customer_username,
  trainer_username,
  product_id,
  sales_username,
  purchase_date,
  total_sessions,
  used_sessions,
  price_paid,
  discount_amount,
  status
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);