-- name: CreateLift :one
INSERT INTO lift (
  exercise_name,
  weight_lifted,
  reps,
  user_id,
  workout_id
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetLift :one
SELECT * FROM lift
WHERE user_id = $1
AND id = $2
LIMIT 1;

-- name: ListLifts :many
SELECT * FROM lift
WHERE user_id = $1
ORDER BY weight_lifted 
LIMIT $2
OFFSET $3;

-- name: ListWeightPRs :many
SELECT * FROM lift
WHERE user_id = $1
ORDER BY weight_lifted DESC
LIMIT $2
OFFSET $3;

-- name: ListWeightPRsByExercise :many
SELECT * FROM lift
WHERE user_id = $1 AND exercise_name = $2
ORDER BY weight DESC
LIMIT $3
OFFSET $4;

-- name: ListWeightPRsByMuscleGroup :many
SELECT l.id, l.exercise_name, weight_lifted, reps, ex.muscle_group FROM lift as l
JOIN exercise AS ex on l.exercise_name = ex.exercise_name 
WHERE ex.muscle_group = $1
AND l.user_id = $2
ORDER BY weight_lifted DESC
LIMIT $3
OFFSET $4;

-- name: ListRepPRs :many
SELECT * FROM lift 
WHERE user_id = $1
ORDER BY reps DESC
LIMIT $2
OFFSET $3;

-- name: UpdateLiftWeight :one
UPDATE lift SET
weight_lifted = $1
WHERE id = $2 AND
user_id = $3
RETURNING *;

-- name: UpdateReps :one
UPDATE lift SET
reps = $1
WHERE id = $2 AND
user_id = $3
RETURNING *;

-- name: DeleteLift :exec
DELETE FROM lift WHERE id = $1;
