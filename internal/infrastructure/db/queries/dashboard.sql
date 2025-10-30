-- name: TotalRevenue :one
SELECT CAST(COALESCE(SUM(amt), 0) AS SIGNED) AS total
FROM (
  SELECT Price_Paid AS amt
  FROM CUSTOMER_DURATIONS
  WHERE Purchase_Date BETWEEN sqlc.arg(start) AND sqlc.arg(end)
  UNION ALL
  SELECT Price_Paid AS amt
  FROM CUSTOMER_SESSIONS
  WHERE Purchase_Date BETWEEN sqlc.arg(start) AND sqlc.arg(end)
) t;

-- name: ActiveMembersToday :one
SELECT CAST(COUNT(DISTINCT Customer_Username) AS SIGNED) AS c
FROM CUSTOMER_DURATIONS
WHERE Status = 'ACTIVE'
  AND CURDATE() BETWEEN DATE(Start_Date) AND DATE(End_Date);

-- name: CheckinsToday :one
SELECT CAST(COUNT(*) AS SIGNED) AS c
FROM CUSTOMER_LOGS
WHERE Log_Type = 'CHECK_IN'
  AND DATE(`Timestamp`) = CURDATE();

-- name: CompletedPTInRange :one
SELECT CAST(COUNT(*) AS SIGNED) AS c
FROM TRAINING_SCHEDULES
WHERE Schedule_Type = 'APPOINTMENT'
  AND Start_Time BETWEEN sqlc.arg(start) AND sqlc.arg(end)
  AND Start_Time < NOW();

-- name: TopSellingProducts :many
SELECT p.Name AS name, COUNT(*) AS units
FROM (
  SELECT Product_Id FROM CUSTOMER_DURATIONS
  WHERE Purchase_Date BETWEEN sqlc.arg(start) AND sqlc.arg(end)
  UNION ALL
  SELECT Product_Id FROM CUSTOMER_SESSIONS
  WHERE Purchase_Date BETWEEN sqlc.arg(start) AND sqlc.arg(end)
) s
JOIN PRODUCTS p ON p.Product_Id = s.Product_Id
GROUP BY p.Name
ORDER BY units DESC
LIMIT 5;

-- name: NewMembersInRange :one
SELECT CAST(COUNT(*) AS SIGNED) AS c
FROM USERS
WHERE Role = 'CUSTOMER'
  AND Created_At BETWEEN sqlc.arg(start) AND sqlc.arg(end);