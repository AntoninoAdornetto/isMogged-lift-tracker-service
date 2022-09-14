-- name: CreateExersise :one
INSERT INTO exersise (
  exersise_name,
  muscle_group
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetExersise :one
SELECT * FROM exersise
WHERE exersise_name = ($1) LIMIT 1;

-- name: GetMuscleGroupExersises :many
SELECT * FROM exersise 
WHERE muscle_group = ($1)
ORDER BY exersise_name;

-- name: UpdateExersiseName :exec
UPDATE exersise SET
exersise_name = ($1);

-- name: UpdateExersiseMuscleGroup :exec
UPDATE exersise SET
muscle_group = ($1);

-- name: DeleteExersise :exec
DELETE FROM exersise WHERE exersise_name = ($1);