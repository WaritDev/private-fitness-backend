-- Payment Slip Verification Queries

-- Q_VERIFY_1: Insert payment verification log
-- name: InsertPaymentVerification :execresult
INSERT INTO payment_verifications (
  customer_username,
  product_id,
  amount,
  slip_file_path,
  verification_status
) VALUES (?, ?, ?, ?, ?);

-- Q_VERIFY_2: Check duplicate payment slip (by customer + product + amount within 24 hours)
-- name: CheckDuplicatePayment :one
SELECT COUNT(id) AS duplicate_count
FROM payment_verifications
WHERE customer_username = ?
  AND product_id = ?
  AND ABS(amount - ?) < 0.01
  AND verification_status = 'VERIFIED'
  AND created_at >= DATE_SUB(NOW(), INTERVAL 24 HOUR);

-- Q_VERIFY_3: Update payment verification status after Slip2Go verification
-- name: UpdatePaymentVerificationStatus :exec
UPDATE payment_verifications
SET 
  verification_status = ?,
  slip_id = ?,
  slip2go_response = ?,
  verified_at = NOW(),
  updated_at = NOW()
WHERE id = ?;

-- Q_VERIFY_4: Get payment verification by ID
-- name: GetPaymentVerificationById :one
SELECT
  id,
  customer_username,
  product_id,
  amount,
  slip_file_path,
  slip_id,
  verification_status,
  slip2go_response,
  verified_at,
  created_at,
  updated_at
FROM payment_verifications
WHERE id = ?;

-- Q_VERIFY_5: List all payment verifications for a customer
-- name: ListPaymentVerificationsByCustomer :many
SELECT
  id,
  customer_username,
  product_id,
  amount,
  slip_file_path,
  slip_id,
  verification_status,
  verified_at,
  created_at
FROM payment_verifications
WHERE customer_username = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- Q_VERIFY_6: Count payment verifications by status
-- name: CountPaymentVerificationsByStatus :one
SELECT COUNT(id) AS total
FROM payment_verifications
WHERE verification_status = ?;
