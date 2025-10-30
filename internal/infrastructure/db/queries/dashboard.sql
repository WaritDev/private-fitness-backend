-- name: ActiveMembersToday :one
SELECT CAST(COUNT(DISTINCT customer_username) AS SIGNED) AS c
FROM customer_durations
WHERE status = 'ACTIVE'
  AND CURDATE() BETWEEN DATE(start_date) AND DATE(end_date);

-- name: CheckinsToday :one
SELECT CAST(COUNT(*) AS SIGNED) AS c
FROM customer_logs
WHERE log_type = 'CHECK_IN'
  AND DATE(`created_at`) = CURDATE();

-- name: CompletedPTInRange :one
SELECT CAST(COUNT(*) AS SIGNED) AS c
FROM training_schedules
WHERE schedule_type = 'APPOINTMENT'
  AND start_time BETWEEN sqlc.arg(start) AND sqlc.arg(end)
  AND start_time < NOW();

-- name: NewMembersInRange :one
SELECT CAST(COUNT(*) AS SIGNED) AS c
FROM users
WHERE role = 'CUSTOMER'
  AND created_at BETWEEN sqlc.arg(start) AND sqlc.arg(end);

-- name: RevenueDurations :one
SELECT CAST(COALESCE(SUM(price_paid), 0) AS SIGNED) AS total
FROM customer_durations
WHERE purchase_date BETWEEN sqlc.arg(start) AND sqlc.arg(end);

-- name: RevenueSessions :one
SELECT CAST(COALESCE(SUM(price_paid), 0) AS SIGNED) AS total
FROM customer_sessions
WHERE purchase_date BETWEEN sqlc.arg(start) AND sqlc.arg(end);

-- name: TopSellingProductsDurations :many
SELECT p.name AS name, COUNT(*) AS units
FROM customer_durations d
JOIN products p ON p.id = d.product_id
WHERE d.purchase_date BETWEEN sqlc.arg(start) AND sqlc.arg(end)
GROUP BY p.name
ORDER BY units DESC;

-- name: TopSellingProductsSessions :many
SELECT p.name AS name, COUNT(*) AS units
FROM customer_sessions s
JOIN products p ON p.id = s.product_id
WHERE s.purchase_date BETWEEN sqlc.arg(start) AND sqlc.arg(end)
GROUP BY p.name
ORDER BY units DESC;