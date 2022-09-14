-- name: CreateAccount :one
INSERT INTO accounts (
  lifter,
  birth_date,
  weight,
  start_date
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY lifter
LIMIT $1
OFFSET $2;

-- name: UpdateAccountWeight :exec
UPDATE accounts SET
weight = $1 WHERE 
id = $2;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;
