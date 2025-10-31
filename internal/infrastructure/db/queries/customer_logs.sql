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

-- name: CreatePendingCheckInLog :exec
-- สร้าง pending check-in log สำหรับลูกค้าที่สแกน QR Code
-- NOTE: ต้องรัน migration ก่อน (เพิ่ม status และ schedule_id columns)
-- หากยังไม่ได้รัน migration ให้ใช้ CreateCustomerLog แทน และจัดการใน use case layer
INSERT INTO customer_logs (
  customer_username,
  log_type,
  status,
  schedule_id
) VALUES (
  ?, 'CHECK_IN', 'PENDING', ?
);

-- name: GetPendingCheckInsByTrainer :many
-- ดึง pending check-ins ของลูกค้าที่มี schedule กับ trainer นี้
SELECT 
  cl.id AS log_id,
  cl.customer_username,
  u.first_name AS customer_first_name,
  u.last_name AS customer_last_name,
  cl.created_at AS checkin_time,
  ts.id AS schedule_id,
  ts.start_time AS appointment_start_time,
  ts.end_time AS appointment_end_time,
  cs.id AS session_id,
  cs.total_sessions,
  cs.used_sessions
FROM customer_logs cl
JOIN training_schedules ts ON ts.id = cl.schedule_id
JOIN users u ON u.username = cl.customer_username
LEFT JOIN customer_sessions cs ON cs.id = ts.session_id
WHERE ts.trainer_username = ?
  AND cl.log_type = 'CHECK_IN'
  AND cl.status = 'PENDING'
  AND DATE(ts.start_time) = CURDATE()
ORDER BY cl.created_at DESC;

-- name: UpdateCheckInLogStatus :execresult
-- อัปเดต status ของ check-in log จาก PENDING เป็น CONFIRMED
UPDATE customer_logs
SET status = 'CONFIRMED'
WHERE id = ? AND status = 'PENDING';