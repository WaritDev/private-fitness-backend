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
  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);