-- name: CreateAccount :one
INSERT INTO accounts (
  name,
  email,
  password,
  weight,
  body_fat
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts
LIMIT $1
OFFSET $2;

-- name: UpdateWeight :exec
UPDATE accounts SET
weight = $1 WHERE 
id = $2 RETURNING *;

-- name: DeleteAccount :one
DELETE FROM accounts WHERE id = $1 RETURNING *;
