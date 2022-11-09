-- name: CreateSet :one
INSERT INTO set (user_id) 
VALUES ($1)
RETURNING *;

-- name: GetSet :one
SELECT * FROM set 
WHERE id = $1 LIMIT 1;

-- name: GetLiftSets :many
SELECT exercise_name, weight, reps, date_lifted, set_id
FROM set JOIN lift ON set.id = lift.set_id WHERE set.id = $1;

-- name: DeleteSet :exec
DELETE FROM set WHERE id = $1;
