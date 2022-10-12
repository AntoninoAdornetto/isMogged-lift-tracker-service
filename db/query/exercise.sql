-- name: CreateExercise :one
INSERT INTO exercise (
  exercise_name,
  muscle_group
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetExercise :one
SELECT * FROM exercise
WHERE exercise_name = ($1) LIMIT 1;

-- name: ListExercises :many
SELECT * FROM exercise
ORDER BY exercise_name
LIMIT $1
OFFSET $2;

-- name: GetMuscleGroupExercises :many
SELECT * FROM exercise 
WHERE muscle_group = ($1)
ORDER BY exercise_name;

-- name: UpdateExerciseName :exec
UPDATE exercise SET
exercise_name = ($1)
WHERE exercise_name = ($2);

-- name: UpdateExerciseMuscleGroup :exec
UPDATE exercise SET
muscle_group = ($1)
WHERE exercise_name = ($2);

-- name: DeleteExercise :exec
DELETE FROM exercise WHERE exercise_name = ($1);
