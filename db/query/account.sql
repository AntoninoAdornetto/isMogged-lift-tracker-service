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

-- name: GetAccountByEmail :one
SELECT id, email, password, password_changed_at, start_date FROM 
accounts WHERE email = $1 LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts
WHERE id = $1
LIMIT $2
OFFSET $3;

-- name: UpdateWeight :exec
UPDATE accounts SET
weight = $1 WHERE 
id = $2 RETURNING *;

-- name: DeleteAccount :one
DELETE FROM accounts WHERE id = $1 RETURNING *;
