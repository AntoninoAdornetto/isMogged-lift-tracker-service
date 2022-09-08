-- name: CreateLift :one
INSERT INTO lift (
  exersise_name,
  weight,
  reps,
  user_id
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetLift :one
SELECT * FROM lift 
WHERE id = $1 LIMIT 1;

-- name: GetLargeWeightLifts :many
SELECT * FROM lift 
WHERE user_id = $1
ORDER BY weight;

-- name: GetHighRepLifts :many
SELECT * FROM lift 
WHERE user_id = $1
ORDER BY reps;

-- name: UpdateWeight :exec
UPDATE lift SET
weight = $1
WHERE id = $2 AND
user_id = $3;

-- name: UpdateReps :exec
UPDATE lift SET
reps = $1
WHERE id = $2 AND
user_id = $3;

-- name: DeleteLift :exec
DELETE FROM lift WHERE id = $1;
