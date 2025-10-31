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

-- name: ListCustomerDurations :many
SELECT
  cd.id,
  cd.customer_username,
  u.first_name  AS customer_first_name,
  u.last_name   AS customer_last_name,
  cd.product_id,
  p.name        AS product_name,
  p.type,
  p.category,
  p.duration_days,
  cd.sales_username,
  cd.purchase_date,
  cd.start_date,
  cd.end_date,
  cd.price_paid,
  cd.discount_amount,
  cd.status
FROM customer_durations cd
JOIN customers c ON c.username = cd.customer_username
JOIN users     u ON u.username = c.username
JOIN products  p ON p.id = cd.product_id
ORDER BY cd.created_at DESC, cd.id DESC
LIMIT ? OFFSET ?;

-- name: CountCustomerDurations :one
SELECT COUNT(cd.id) AS total_items
FROM customer_durations cd
JOIN customers c ON c.username = cd.customer_username
JOIN users     u ON u.username = c.username
JOIN products  p ON p.id = cd.product_id;

-- name: GetCustomerDurationByID :one
SELECT
  cd.id,
  cd.customer_username,
  cd.product_id,
  cd.sales_username,
  cd.purchase_date,
  cd.start_date,
  cd.end_date,
  cd.price_paid,
  cd.discount_amount,
  cd.status
FROM customer_durations cd
WHERE cd.id = ?;

-- name: GetDurationDaysForDurationID :one
SELECT p.duration_days
FROM customer_durations cd
JOIN products p ON p.id = cd.product_id
WHERE cd.id = ?
  AND p.type = 'DURATION'
  AND p.is_active = 1
  AND p.duration_days IS NOT NULL
LIMIT 1;

-- name: UpdateCustomerDurationEditableFields :exec
UPDATE customer_durations cd
JOIN products p ON p.id = cd.product_id
SET
  cd.start_date      = STR_TO_DATE(?, '%Y-%m-%d'),
  cd.end_date        = DATE_ADD(STR_TO_DATE(?, '%Y-%m-%d'), INTERVAL (p.duration_days - 1) DAY),
  cd.price_paid      = ?,
  cd.discount_amount = ?,
  cd.status          = ?,
  cd.updated_at      = NOW()
WHERE cd.id = ?;

-- name: DeleteCustomerDurationByID :execresult
DELETE FROM customer_durations
WHERE id = ?;


-- name: GetCustomerActiveDuration :many
-- ดึงข้อมูล Duration packages ที่ยัง ACTIVE ของลูกค้า
-- JOIN กับ PRODUCTS เพื่อดึงชื่อแพ็กเกจและคำนวณ วันคงเหลือที่ยังใช้งานได้
SELECT 
  cd.id,
  cd.customer_username,
  cd.product_id,
  p.name AS product_name,
  p.duration_days,
  cd.sales_username,
  cd.purchase_date,
  cd.start_date,
  cd.end_date,
  DATEDIFF(cd.end_date, CURDATE()) AS days_remaining,
  cd.price_paid,
  cd.discount_amount,
  cd.status,
  cd.created_at
FROM customer_durations cd
JOIN products p ON cd.product_id = p.id
WHERE cd.customer_username = ?
  AND cd.status = 'ACTIVE'
ORDER BY cd.created_at DESC;


-- name: RenewCustomerDuration :exec
-- Customer Self-Purchase: ลูกค้าซื้อแพ็กเกจ Duration เพิ่มเอง
-- INSERT แพ็กเกจใหม่ โดย sales_username = NULL, discount_amount = 0
-- Backend คำนวณ start_date = NOW() และ end_date = NOW() + duration_days (ดึงจาก products)
INSERT INTO customer_durations (
  customer_username,
  product_id,
  sales_username,      -- NULL สำหรับการซื้อเอง
  purchase_date,       -- NOW()
  start_date,          -- NOW() (Backend กำหนดให้)
  end_date,            -- NOW() + duration_days (ดึงจาก products.duration_days)
  price_paid,          -- list_price (ราคาเต็ม)
  discount_amount,     -- 0 (ไม่มีส่วนลด)
  status               -- 'ACTIVE'
)
SELECT 
  ?,                   -- customer_username
  ?,                   -- product_id
  NULL,                -- sales_username (self-purchase)
  NOW(),               -- purchase_date
  NOW(),               -- start_date
  DATE_ADD(NOW(), INTERVAL p.duration_days DAY), -- end_date (NOW + duration_days)
  ?,                   -- price_paid (list_price)
  '0.00',              -- discount_amount
  'ACTIVE'             -- status
FROM products p
WHERE p.id = ?         -- product_id (ตรวจสอบและดึง duration_days)
  AND p.type = 'DURATION'
  AND p.is_active = 1
  AND p.duration_days IS NOT NULL
LIMIT 1;

