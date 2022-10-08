-- name: CreateMuscleGroup :one
INSERT INTO muscle_groups (
  group_name
) VALUES (
  $1
)
RETURNING *;

-- name: GetMuscleGroup :one
SELECT * FROM muscle_groups
WHERE group_name = $1;

-- name: GetMuscleGroups :many
SELECT * FROM muscle_groups
ORDER BY group_name;

-- name: UpdateGroup :one
UPDATE muscle_groups SET group_name = $1 WHERE group_name = $2 RETURNING *;

-- name: DeleteGroup :exec
DELETE FROM muscle_groups WHERE group_name = $1;
