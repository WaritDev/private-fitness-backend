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
