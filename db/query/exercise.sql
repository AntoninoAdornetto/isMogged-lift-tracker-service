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

-- name: UpdateExercise :one
UPDATE exercise SET
name = COALESCE(NULLIF($1, ''), name),
muscle_group = COALESCE(NULLIF($2, ''), muscle_group),
category = COALESCE(NULLIF($3, ''), category)
WHERE name = $4
RETURNING *;

-- name: DeleteExercise :exec
DELETE FROM exercise WHERE name = ($1);
