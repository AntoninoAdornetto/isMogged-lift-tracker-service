-- name: CreateMuscleGroup :one
INSERT INTO muscle_group (
  name
) VALUES (
  $1
)
RETURNING *;

-- name: GetMuscleGroup :one
SELECT * FROM muscle_group
WHERE name = $1;

-- name: GetMuscleGroups :many
SELECT * FROM muscle_group
ORDER BY name;

-- name: UpdateGroup :one
UPDATE muscle_group SET name = $1 WHERE name = $2 RETURNING *;

-- name: DeleteGroup :one
DELETE FROM muscle_group WHERE name = $1 RETURNING *;
