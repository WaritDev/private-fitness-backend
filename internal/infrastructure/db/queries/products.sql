-- name: ListDurations :many
SELECT
  id,             
  name,             
  type,            
  category,         
  list_price,       
  duration_days,    
  session_amount,   
  is_active,
  payment_account_id
FROM products
WHERE
  type = 'DURATION'
  AND is_active = TRUE
ORDER BY list_price ASC;

-- name: ListSessions :many
SELECT
  id,
  name,
  type,
  category,
  list_price,
  duration_days,
  session_amount,
  is_active,
  payment_account_id
FROM products
WHERE
  type = 'SESSION'
  AND is_active = TRUE
ORDER BY list_price ASC;

-- name: ListAllProducts :many
SELECT
  id,
  name,
  type,
  category,
  list_price,
  duration_days,
  session_amount,
  is_active,
  payment_account_id,
  created_at,
  updated_at
FROM products
WHERE is_active = TRUE
ORDER BY category, list_price ASC;

-- name: GetProductById :one
SELECT
  id,
  name,
  type,
  category,
  list_price,
  duration_days,
  session_amount,
  is_active,
  payment_account_id,
  created_at,
  updated_at
FROM products
WHERE id = ?
LIMIT 1;
