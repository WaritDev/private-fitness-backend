-- name: CreateTrainingSchedule :exec
INSERT INTO training_schedules (
  trainer_username,
  customer_username,
  session_id,
  start_time,
  end_time,
  schedule_type
) VALUES (
  ?, ?, ?, ?, ?, ?
);

-- name: GetDayOffSchedules :many
-- Q3C.3a - ดึงวันหยุดหรือช่วงเวลาที่ไม่รับนัด (DAY_OFF)
SELECT
  start_time,
  end_time
FROM training_schedules
WHERE trainer_username = ?
  AND schedule_type = 'DAY_OFF'
  AND start_time < ?
  AND end_time > ?;

-- name: GetAppointmentSchedules :many
-- Q3C.3b - ดึงนัดที่ถูกจองแล้ว (APPOINTMENT)
SELECT
  id,
  start_time,
  end_time,
  customer_username
FROM training_schedules
WHERE trainer_username = ?
  AND schedule_type = 'APPOINTMENT'
  AND start_time < ?
  AND end_time > ?;

-- name: CheckTimeSlotAvailability :one
-- Q2C.6 - ตรวจสอบว่าช่วงเวลาที่เลือกยังว่างอยู่จริง
-- คืนค่า overlapped_count ถ้าเป็น 0 แสดงว่ายังว่างอยู่
-- Logic: ช่วงเวลาซ้อนทับกันเมื่อ start_time < endTime AND end_time > startTime
SELECT COUNT(id) AS overlapped_count
FROM training_schedules
WHERE trainer_username = ?
  AND start_time < ?
  AND end_time > ?;

-- name: BookAppointment :exec
-- Q2C.7 - บันทึกการจองนัด (APPOINTMENT)
INSERT INTO training_schedules (
  trainer_username,
  customer_username,
  session_id,
  start_time,
  end_time,
  schedule_type
) VALUES (
  ?, ?, ?, ?, ?, 'APPOINTMENT'
);

-- name: GetAppointmentById :one
-- ดึงข้อมูลการจองตาม ID เพื่อตรวจสอบก่อนยกเลิก
SELECT
  id,
  trainer_username,
  customer_username,
  session_id,
  start_time,
  end_time,
  schedule_type
FROM training_schedules
WHERE id = ?
  AND schedule_type = 'APPOINTMENT';

-- name: DeleteAppointment :exec
-- ยกเลิกการจอง - ลบ training_schedule record
-- Logic: DELETE FROM training_schedules WHERE id = ? AND schedule_type = 'APPOINTMENT'
-- Note: ใช้ id เพียงอย่างเดียว เพราะ id เป็น Primary Key (unique)
-- การตรวจสอบ customer_username ทำใน use case layer ก่อนเรียก query นี้
DELETE FROM training_schedules
WHERE id = ?
  AND schedule_type = 'APPOINTMENT';

-- Use Case 3P: Manage Day-Offs

-- Q3P.1: Get all day-offs for a trainer
-- name: GetTrainerDayOffs :many
SELECT
  id,
  trainer_username,
  start_time,
  end_time
FROM training_schedules
WHERE trainer_username = ?
  AND schedule_type = 'DAY_OFF'
ORDER BY start_time DESC;

-- Q3P.2: Check if day-off already exists for the same date
-- name: CheckDayOffDuplicate :one
SELECT COUNT(id) AS duplicate_count
FROM training_schedules
WHERE trainer_username = ?
  AND schedule_type = 'DAY_OFF'
  AND DATE(start_time) = ?;

-- Q6C.5: Check if day-off overlaps with existing appointments
-- name: CheckDayOffAppointmentOverlap :one
SELECT COUNT(id) AS overlapped_count
FROM training_schedules
WHERE trainer_username = ?
  AND schedule_type = 'APPOINTMENT'
  AND start_time < ?
  AND end_time > ?;

-- Q3P.4: Create day-off
-- name: CreateDayOff :exec
INSERT INTO training_schedules (
  trainer_username,
  start_time,
  end_time,
  schedule_type
) VALUES (
  ?, ?, ?, 'DAY_OFF'
);

-- Q3P.5: Delete day-off
-- name: DeleteDayOff :exec
DELETE FROM training_schedules
WHERE id = ?
  AND schedule_type = 'DAY_OFF';

-- Find customer's schedule for today (for check-in flow)
-- name: GetCustomerScheduleForToday :one
SELECT 
  id,
  trainer_username,
  customer_username,
  session_id,
  start_time,
  end_time
FROM training_schedules
WHERE customer_username = ?
  AND schedule_type = 'APPOINTMENT'
  AND DATE(start_time) = CURDATE()
  AND start_time <= NOW()
  AND end_time >= NOW()
ORDER BY start_time DESC
LIMIT 1;

-- Get trainer's appointments with pending check-ins (for trainer calendar)
-- name: GetTrainerAppointmentsWithPendingCheckIns :many
-- Note: Requires customer_logs.status and schedule_id columns (run migration first)
SELECT 
  ts.id AS schedule_id,
  ts.customer_username,
  u.first_name AS customer_first_name,
  u.last_name AS customer_last_name,
  ts.start_time,
  ts.end_time,
  ts.session_id,
  cs.total_sessions,
  cs.used_sessions,
  COALESCE(cl.status, 'NONE') AS checkin_status,
  cl.id AS checkin_log_id,
  cl.created_at AS checkin_time
FROM training_schedules ts
JOIN users u ON u.username = ts.customer_username
LEFT JOIN customer_sessions cs ON cs.id = ts.session_id
LEFT JOIN customer_logs cl ON cl.schedule_id = ts.id 
  AND cl.log_type = 'CHECK_IN' 
  AND cl.status = 'PENDING'
WHERE ts.trainer_username = ?
  AND ts.schedule_type = 'APPOINTMENT'
  AND DATE(ts.start_time) >= CURDATE()
ORDER BY ts.start_time ASC;
