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

-- name: CreateLifts :many
INSERT INTO lift (
  exercise_name,
  weight_lifted,
  reps,
  user_id,
  workout_id
) VALUES (
  UNNEST(@exerciseNames::VARCHAR[]),
  UNNEST(@weights::REAL[]),
  UNNEST(@reps::SMALLINT[]),
  UNNEST(@user_id::UUID[]),
  UNNEST(@workout_id::UUID[])
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
ORDER BY exercise_name 
LIMIT $2
OFFSET $3;

-- name: ListPRs :many
SELECT * FROM lift
WHERE user_id = $1
ORDER BY
  CASE
    WHEN $2 = 'weight' THEN weight_lifted
    WHEN $3 = 'reps' THEN reps
END DESC
LIMIT $4
OFFSET $5;

-- name: ListPRsByExercise :many
SELECT * FROM lift
WHERE user_id = $1 AND exercise_name = $2
ORDER BY
  CASE
    WHEN $3 = 'weight' THEN weight_lifted
    WHEN $4 = 'reps' THEN reps 
END DESC
LIMIT $5
OFFSET $6;

-- name: ListPRsByMuscleGroup :many
SELECT l.id, l.exercise_name, weight_lifted, reps, ex.muscle_group FROM lift AS l
JOIN exercise AS ex on l.exercise_name = ex.name
WHERE ex.muscle_group = $1
AND l.user_id = $2
ORDER BY
  CASE
    WHEN $3 = 'weight' THEN weight_lifted
    WHEN $4 = 'reps' THEN reps
END DESC
LIMIT $5
OFFSET $6;

-- name: UpdateLift :one
UPDATE lift SET
weight_lifted = COALESCE(NULLIF($1, 0::REAL), weight_lifted),
reps = COALESCE(NULLIF($2, 0::SMALLINT), reps)
WHERE id = $3
RETURNING *;

-- name: DeleteLift :exec
DELETE FROM lift WHERE id = $1;
