-- name: CreatePaymentAccount :exec
INSERT INTO payment_accounts (
  account_name,
  account_number,
  bank_name,
  qr_code_image_url,
  is_active
) VALUES (
  ?,
  ?,
  ?,
  ?,
  ?
);

-- Q5S.1: ดึงข้อมูลสินค้าและบัญชีรับชำระเงินเพื่อแสดงหน้าชำระเงิน
-- name: GetPaymentInfoByProductId :one
SELECT 
  p.id AS product_id,
  p.name AS product_name,
  p.type AS product_type,
  p.category AS product_category,
  p.list_price,
  p.duration_days,
  p.session_amount,
  pa.id AS payment_account_id,
  pa.account_name,
  pa.account_number,
  pa.bank_name,
  pa.qr_code_image_url,
  pa.is_active AS account_active
FROM products p
JOIN payment_accounts pa ON pa.id = p.payment_account_id
WHERE p.id = ?
  AND p.is_active = TRUE
  AND pa.is_active = TRUE;