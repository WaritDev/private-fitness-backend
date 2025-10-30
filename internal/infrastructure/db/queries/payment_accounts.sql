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