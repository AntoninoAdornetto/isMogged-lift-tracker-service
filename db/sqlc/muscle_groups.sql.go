// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: muscle_groups.sql

package db

import (
	"context"
)

const createMuscleGroup = `-- name: CreateMuscleGroup :one
INSERT INTO muscle_groups (
  group_name
) VALUES (
  $1
)
RETURNING id, group_name
`

func (q *Queries) CreateMuscleGroup(ctx context.Context, groupName string) (MuscleGroup, error) {
	row := q.db.QueryRowContext(ctx, createMuscleGroup, groupName)
	var i MuscleGroup
	err := row.Scan(&i.ID, &i.GroupName)
	return i, err
}

const deleteGroup = `-- name: DeleteGroup :exec
DELETE FROM muscle_groups WHERE group_name = $1
`

func (q *Queries) DeleteGroup(ctx context.Context, groupName string) error {
	_, err := q.db.ExecContext(ctx, deleteGroup, groupName)
	return err
}

const getMuscleGroup = `-- name: GetMuscleGroup :one
SELECT id, group_name FROM muscle_groups
WHERE group_name = $1
`

func (q *Queries) GetMuscleGroup(ctx context.Context, groupName string) (MuscleGroup, error) {
	row := q.db.QueryRowContext(ctx, getMuscleGroup, groupName)
	var i MuscleGroup
	err := row.Scan(&i.ID, &i.GroupName)
	return i, err
}

const getMuscleGroups = `-- name: GetMuscleGroups :many
SELECT id, group_name FROM muscle_groups
ORDER BY group_name
`

func (q *Queries) GetMuscleGroups(ctx context.Context) ([]MuscleGroup, error) {
	rows, err := q.db.QueryContext(ctx, getMuscleGroups)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MuscleGroup
	for rows.Next() {
		var i MuscleGroup
		if err := rows.Scan(&i.ID, &i.GroupName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateGroup = `-- name: UpdateGroup :exec
UPDATE muscle_groups SET group_name = $1
`

func (q *Queries) UpdateGroup(ctx context.Context, groupName string) error {
	_, err := q.db.ExecContext(ctx, updateGroup, groupName)
	return err
}
