-- name: CreateWorkout :one
INSERT INTO workout (user_id, start_time) 
VALUES ($1, $2)
RETURNING *;

-- name: GetWorkout :many
SELECT w.id, exercise_name, weight_lifted, reps, start_time, finish_time, l.user_id
FROM workout AS w
JOIN lift AS l ON l.workout_id = w.id
WHERE w.id = $1;

-- name: UpdateFinishTime :one
UPDATE workout SET
finish_time = $1
WHERE id = $2
RETURNING *;

-- name: ListWorkouts :many
SELECT * FROM workout
WHERE user_id = $1
ORDER BY start_time DESC
LIMIT $2
OFFSET $3;

-- name: DeleteWorkout :exec
DELETE FROM workout
WHERE id = $1;

