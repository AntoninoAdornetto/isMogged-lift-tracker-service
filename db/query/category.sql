-- name: CreateCategory :one
INSERT INTO category (
  name
) VALUES ($1)
RETURNING *;

-- name: GetCategory :one
SELECT * FROM category
WHERE id = $1 LIMIT 1;

-- name: ListCategories :many
SELECT * FROM category;

-- name: UpdateCategory :exec
UPDATE category SET
name = $1 WHERE 
id = $2 RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM category WHERE id = $1;
