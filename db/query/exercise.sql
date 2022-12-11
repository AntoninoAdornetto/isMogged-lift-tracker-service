-- name: CreateExercise :one
INSERT INTO exercise (
  name,
  muscle_group,
  category
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetExercise :one
SELECT * FROM exercise
WHERE name = ($1) LIMIT 1;

-- name: ListExercises :many
SELECT * FROM exercise
ORDER BY name 
LIMIT $1
OFFSET $2;

-- name: ListByMuscleGroup :many
SELECT * FROM exercise 
WHERE muscle_group = ($1)
ORDER BY name
LIMIT $2
OFFSET $3;

-- name: UpdateExerciseName :exec
UPDATE exercise SET
name = ($1)
WHERE name = ($2) RETURNING *;

-- name: UpdateMuscleGroup :exec
UPDATE exercise SET
muscle_group = ($1)
WHERE name = ($2) RETURNING *;

-- name: DeleteExercise :exec
DELETE FROM exercise WHERE name = ($1);
