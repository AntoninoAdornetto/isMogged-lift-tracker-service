-- name: CreateLift :one
INSERT INTO lift (
  exersise_name,
  weight,
  reps,
  user_id,
  set_id
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetLift :one
SELECT * FROM lift 
WHERE id = $1 LIMIT 1;

-- name: GetWeightPRs :many
SELECT * FROM lift 
WHERE user_id = $1
ORDER BY weight;

-- name: GetRepPRs :many
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
