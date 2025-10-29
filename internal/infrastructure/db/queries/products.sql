-- name: ListDurations :many
SELECT
  Product_Id,
  Name,
  Product_Type,
  Product_Category,
  List_Price,
  Duration_Days,
  Session_Amount,
  Is_Active
FROM
  PRODUCTS
WHERE
  Product_Type = 'DURATION'
  AND Is_Active = 1
ORDER BY
  List_Price ASC


-- name: ListSessions :many
SELECT
  Product_Id,
  Name,
  Product_Type,
  Product_Category,
  List_Price,
  Duration_Days,
  Session_Amount,
  Is_Active
FROM
  PRODUCTS
WHERE
  Product_Type = 'SESSION'
  AND Is_Active = 1
ORDER BY
  List_Price ASC

