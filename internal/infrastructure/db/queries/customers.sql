-- name: CreateCustomer :exec
INSERT INTO
  customers (
    username,
    health_info,
    address,
    company_name,
    company_position,
    marital_status,
    emergency_contact_name,
    emergency_contact_relationship,
    emergency_contact_phone,
    marketing_source
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: ListCustomers :many
SELECT
  u.username,
  u.first_name,
  u.last_name,
  u.gender,
  u.date_of_birth,
  u.phone_number,
  u.gmail,
  u.is_active,
  c.health_info,
  c.address,
  c.company_name,
  c.company_position,
  c.marital_status,
  c.emergency_contact_name,
  c.emergency_contact_relationship,
  c.emergency_contact_phone,
  c.marketing_source
FROM users u
JOIN customers c ON c.username = u.username
WHERE u.role = 'CUSTOMER'
ORDER BY u.is_active DESC, u.first_name ASC, u.last_name ASC
LIMIT ? OFFSET ?;

-- name: CountCustomers :one
SELECT COUNT(u.username) AS total_items
FROM users u
JOIN customers c ON c.username = u.username
WHERE u.role = 'CUSTOMER';

-- name: GetCustomerByUsername :one
SELECT
  u.username,
  u.first_name,
  u.last_name,
  u.gender,
  u.date_of_birth,
  u.phone_number,
  u.gmail,
  u.is_active,
  c.health_info,
  c.address,
  c.company_name,
  c.company_position,
  c.marital_status,
  c.emergency_contact_name,
  c.emergency_contact_relationship,
  c.emergency_contact_phone,
  c.marketing_source
FROM users u
JOIN customers c ON c.username = u.username
WHERE u.username = ?;

-- name: CheckCustomerPhoneExistsExcept :one
SELECT COUNT(u.username) AS count
FROM users u
WHERE u.phone_number = ?
  AND u.username <> ?;

-- name: CheckCustomerGmailExistsExcept :one
SELECT COUNT(u.username) AS count
FROM users u
WHERE LOWER(u.gmail) = LOWER(?)
  AND u.username <> ?;

-- อัปเดต users: ไม่มีการเปลี่ยน password
-- name: UpdateCustomerUserNoPassword :exec
UPDATE users
SET first_name    = ?,
    last_name     = ?,
    gender        = ?,
    date_of_birth = ?,
    phone_number  = ?,
    gmail         = ?,
    is_active     = ?,
    updated_at    = CURRENT_TIMESTAMP
WHERE username = ?;

-- อัปเดต users: มีการรีเซ็ตรหัสผ่าน
-- name: UpdateCustomerUserWithPassword :exec
UPDATE users
SET password      = ?,
    first_name    = ?,
    last_name     = ?,
    gender        = ?,
    date_of_birth = ?,
    phone_number  = ?,
    gmail         = ?,
    is_active     = ?,
    updated_at    = CURRENT_TIMESTAMP
WHERE username = ?;

-- อัปเดตตาราง customers
-- name: UpdateCustomersDetail :exec
UPDATE customers
SET health_info                    = ?,
    address                        = ?,
    company_name                   = ?,
    company_position               = ?,
    marital_status                 = ?,
    emergency_contact_name         = ?,
    emergency_contact_relationship = ?,
    emergency_contact_phone        = ?,
    marketing_source               = ?
WHERE username = ?;

-- name: DeleteCustomerLogsByCustomer :exec
DELETE FROM customer_logs
WHERE customer_username = ?;

-- name: DeleteTrainingScheduleByCustomer :exec
DELETE FROM training_schedules
WHERE customer_username = ?;

-- name: DeleteCustomerSessionByCustomer :exec
DELETE FROM customer_sessions
WHERE customer_username = ?;

-- name: DeleteCustomerDurationByCustomer :exec
DELETE FROM customer_durations
WHERE customer_username = ?;

-- name: DeleteCustomerByUsername :exec
DELETE FROM customers
WHERE username = ?;