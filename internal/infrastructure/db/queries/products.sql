-- name: ListDurations :many
SELECT
  id,             
  name,             
  type,            
  category,         
  list_price,       
  duration_days,    
  session_amount,   
  is_active
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
  is_active
FROM products
WHERE
  type = 'SESSION'
  AND is_active = TRUE
ORDER BY list_price ASC;

