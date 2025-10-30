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

-- name: CheckBookingPermission :one
-- ตรวจสอบสิทธิ์การจองของ Customer
-- ต้องมี Session package แบบ ACTIVE และยังมีสิทธิ์คงเหลือ (used_sessions < total_sessions)
SELECT COUNT(*) as has_permission
FROM customer_sessions
WHERE customer_username = ?
  AND status = 'ACTIVE'
  AND used_sessions < total_sessions;

-- name: IncrementUsedSessions :exec
-- Q3C.6 - อัปเดตจำนวนครั้งที่ใช้ไปแล้ว
UPDATE customer_sessions
SET used_sessions = used_sessions + 1
WHERE id = ?;

-- name: GetActiveSessionByCustomer :one
-- หา Session package ACTIVE ของ Customer (สำหรับการจอง)
SELECT
  id,
  customer_username,
  trainer_username,
  total_sessions,
  used_sessions
FROM customer_sessions
WHERE customer_username = ?
  AND status = 'ACTIVE'
  AND used_sessions < total_sessions
ORDER BY created_at DESC
LIMIT 1;

-- name: DecrementUsedSessions :exec
-- ยกเลิกการจอง - ลดจำนวนครั้งที่ใช้ไป (used_sessions - 1)
-- Logic: UPDATE customer_sessions SET used_sessions = used_sessions - 1 WHERE id = ? AND used_sessions > 0
-- Note: id คือ customer_sessions.id (ได้จาก training_schedules.session_id)
-- ต้องมี used_sessions > 0 เพื่อป้องกันค่าติดลบ
UPDATE customer_sessions
SET used_sessions = used_sessions - 1
WHERE id = ?
  AND used_sessions > 0;

-- name: GetCustomerActiveSessions :many
-- ดึงข้อมูล Session packages ที่ยัง ACTIVE ของลูกค้า
-- JOIN กับ PRODUCTS เพื่อดึงชื่อแพ็กเกจและคำนวณ sessions คงเหลือ
SELECT 
  cs.id,
  cs.customer_username,
  cs.trainer_username,
  cs.product_id,
  p.name AS product_name,
  cs.total_sessions,
  cs.used_sessions,
  (cs.total_sessions - cs.used_sessions) AS sessions_remaining,
  cs.purchase_date,
  cs.price_paid,
  cs.discount_amount,
  cs.status,
  cs.created_at
FROM customer_sessions cs
JOIN products p ON cs.product_id = p.id
WHERE cs.customer_username = ?
  AND cs.status = 'ACTIVE'
ORDER BY cs.created_at DESC;

-- name: ListCustomerSessions :many
SELECT
  cs.id,
  cs.customer_username,
  cu.first_name  AS customer_first_name,
  cu.last_name   AS customer_last_name,
  cs.trainer_username,
  tu.first_name  AS trainer_first_name,
  tu.last_name   AS trainer_last_name,
  cs.product_id,
  p.name         AS product_name,
  p.type,
  p.category,
  p.session_amount,
  cs.sales_username,
  cs.purchase_date,
  cs.total_sessions,
  cs.used_sessions,
  (cs.total_sessions - cs.used_sessions) AS remaining_sessions,
  cs.price_paid,
  cs.discount_amount,
  cs.status
FROM customer_sessions cs
JOIN users cu ON cu.username = cs.customer_username
JOIN users tu ON tu.username = cs.trainer_username
JOIN products p ON p.id = cs.product_id
ORDER BY cs.created_at DESC, cs.id DESC
LIMIT ? OFFSET ?;

-- name: CountCustomerSessions :one
SELECT COUNT(cs.id) AS total_items
FROM customer_sessions cs
JOIN users cu ON cu.username = cs.customer_username
JOIN users tu ON tu.username = cs.trainer_username
JOIN products p ON p.id = cs.product_id;

-- name: GetCustomerSessionByID :one
SELECT
  cs.id,
  cs.customer_username,
  cs.trainer_username,
  cs.product_id,
  cs.sales_username,
  cs.purchase_date,
  cs.total_sessions,
  cs.used_sessions,
  cs.price_paid,
  cs.discount_amount,
  cs.status
FROM customer_sessions cs
WHERE cs.id = ?;

-- name: CheckTrainerExists :one
SELECT COUNT(username) AS cnt
FROM users
WHERE username = ? AND role = 'TRAINER';

-- name: UpdateCustomerSessionEditableFields :exec
UPDATE customer_sessions
SET trainer_username = ?,
    price_paid       = ?,
    discount_amount  = ?,
    status           = ?,
    updated_at       = NOW()
WHERE id = ?;

-- name: DeleteCustomerSessionByID :execresult
DELETE FROM customer_sessions
WHERE id = ?;