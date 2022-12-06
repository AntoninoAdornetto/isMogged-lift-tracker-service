-- name: CreateWorkout :one
INSERT INTO workout (user_id) 
VALUES ($1)
RETURNING *;

-- name: GetWorkout :many
SELECT w.id, exercise_name, weight_lifted, reps, start_time, finish_time, l.user_id
FROM workout AS w
JOIN lift AS l ON w.id = $1 AND l.user_id = $2;

-- name: UpdateDurationEnd :one
UPDATE workout SET
finish_time = $1
WHERE id = $2
RETURNING *;

-- name: ListWorkouts :many
SELECT
w.id, exercise_name, weight_lifted, reps, start_time, finish_time, l.user_id
FROM workout AS w
JOIN lift AS l ON w.id = $1 AND l.user_id = $2
ORDER BY start_time
LIMIT $3
OFFSET $4;

--@todo
--ListWeightPRWorkouts :many
--ListRepPRWorkouts :many

-- name: DeleteWorkout :exec
DELETE FROM workout
WHERE id = $1;

