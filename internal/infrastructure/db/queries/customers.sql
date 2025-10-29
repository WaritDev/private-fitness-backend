-- name: CreateCustomer :one
INSERT INTO
  "CUSTOMER" (
    "Username",
    "Health_Info",
    "Address",
    "Company_Name",
    "Company_Position",
    "Marital_Status",
    "Emergency_Contact_Name",
    "Emergency_Contact_Relationship",
    "Emergency_Contact_Phone",
    "Marketing_Source"
  )
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);