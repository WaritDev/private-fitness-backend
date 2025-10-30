-- name: FindAvailableTrainers :many
-- Q4S.1 แก้ไข: หาเทรนเนอร์ที่ว่างในวันและเวลาที่กำหนด
SELECT 
  u.username,
  u.first_name,
  u.last_name,
  u.created_at
FROM users u
JOIN training_availabilities ta 
  ON ta.trainer_username = u.username
WHERE 
  u.role = 'TRAINER'
  AND u.is_active = TRUE
  AND ta.day_of_week = ?
  AND TIME(?) >= ta.start_time
  AND TIME(?) <= ta.end_time
ORDER BY u.created_at ASC;

-- name: CountAppointmentsOnDate :one
-- นับจำนวนนัดหมายของเทรนเนอร์ในวันที่กำหนด
SELECT COUNT(*) as count
FROM training_schedules
WHERE trainer_username = ?
  AND schedule_type = 'APPOINTMENT'
  AND DATE(start_time) = DATE(?);

-- name: CheckScheduleOverlap :one
-- ตรวจสอบว่ามีนัดซ้อนทับหรือไม่
SELECT COUNT(*) as count
FROM training_schedules
WHERE trainer_username = ?
  AND schedule_type = 'APPOINTMENT'
  AND start_time < ?
  AND end_time > ?;

-- name: ListAllTrainers :many
-- ดึงรายชื่อเทรนเนอร์ทั้งหมดที่ active (สำหรับ dropdown)
SELECT 
  username,
  first_name,
  last_name
FROM users
WHERE role = 'TRAINER'
  AND is_active = TRUE
ORDER BY first_name, last_name, username;
