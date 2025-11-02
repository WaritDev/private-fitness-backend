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


-- name: ListProducts :many
SELECT
  p.id,
  p.name,
  p.type,
  p.category,
  p.list_price,
  p.duration_days,
  p.session_amount,
  p.is_active,
  p.payment_account_id,
  p.created_at,
  p.updated_at,
  (p.is_active = 1) AS is_active_bool
FROM products p
ORDER BY p.is_active DESC, p.type ASC, p.category ASC, p.list_price ASC;

-- name: CheckPaymentAccountActive :one
SELECT COUNT(pa.id) AS account_count
FROM payment_accounts pa
WHERE pa.id = ?
  AND pa.is_active = 1;

-- name: InsertProductDuration :execresult
INSERT INTO products (
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
) VALUES (
  ?,               -- name
  'DURATION',      -- type
  ?,               -- category
  ?,               -- list_price (DECIMAL -> ส่ง string)
  ?,               -- duration_days
  NULL,            -- session_amount
  COALESCE(?, 1),  -- is_active (NULL -> TRUE)
  ?,               -- payment_account_id
  NOW(), NOW()
);

-- name: InsertProductSession :execresult
INSERT INTO products (
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
) VALUES (
  ?,               -- name
  'SESSION',       -- type
  ?,               -- category
  ?,               -- list_price (DECIMAL -> ส่ง string)
  NULL,            -- duration_days
  ?,               -- session_amount
  COALESCE(?, 1),  -- is_active (NULL -> TRUE)
  ?,               -- payment_account_id
  NOW(), NOW()
);

-- name: UpdateProductDuration :exec
UPDATE products
SET name = ?,
    category = ?,
    list_price = ?,
    duration_days = ?,
    session_amount = NULL,
    is_active = ?,
    payment_account_id = ?,
    updated_at = NOW()
WHERE id = ?
  AND type = 'DURATION';

-- name: UpdateProductSession :exec
UPDATE products
SET name = ?,
    category = ?,
    list_price = ?,
    duration_days = NULL,
    session_amount = ?,
    is_active = ?,
    payment_account_id = ?,
    updated_at = NOW()
WHERE id = ?
  AND type = 'SESSION';

-- name: CountProductReferences :one
WITH ref AS (
  SELECT
    (SELECT COUNT(product_id) FROM customer_durations cd WHERE cd.product_id = ?) AS cd_refs,
    (SELECT COUNT(product_id) FROM customer_sessions cs WHERE cs.product_id = ?) AS cs_refs
)
SELECT
  cd_refs AS ref_in_customer_duration,
  cs_refs AS ref_in_customer_session,
  (cd_refs + cs_refs) AS total_refs
FROM ref;

-- name: DeleteProductByID :execresult
DELETE FROM products
WHERE id = ?;