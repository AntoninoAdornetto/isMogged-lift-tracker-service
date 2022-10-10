// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: exersise.sql

package db

import (
	"context"
)

const createExersise = `-- name: CreateExersise :one
INSERT INTO exersise (
  exersise_name,
  muscle_group
) VALUES (
  $1, $2
) RETURNING id, exersise_name, muscle_group
`

type CreateExersiseParams struct {
	ExersiseName string `json:"exersise_name"`
	MuscleGroup  string `json:"muscle_group"`
}

func (q *Queries) CreateExersise(ctx context.Context, arg CreateExersiseParams) (Exersise, error) {
	row := q.db.QueryRowContext(ctx, createExersise, arg.ExersiseName, arg.MuscleGroup)
	var i Exersise
	err := row.Scan(&i.ID, &i.ExersiseName, &i.MuscleGroup)
	return i, err
}

const deleteExersise = `-- name: DeleteExersise :exec
DELETE FROM exersise WHERE exersise_name = ($1)
`

func (q *Queries) DeleteExersise(ctx context.Context, exersiseName string) error {
	_, err := q.db.ExecContext(ctx, deleteExersise, exersiseName)
	return err
}

const getExersise = `-- name: GetExersise :one
SELECT id, exersise_name, muscle_group FROM exersise
WHERE exersise_name = ($1) LIMIT 1
`

func (q *Queries) GetExersise(ctx context.Context, exersiseName string) (Exersise, error) {
	row := q.db.QueryRowContext(ctx, getExersise, exersiseName)
	var i Exersise
	err := row.Scan(&i.ID, &i.ExersiseName, &i.MuscleGroup)
	return i, err
}

const getMuscleGroupExersises = `-- name: GetMuscleGroupExersises :many
SELECT id, exersise_name, muscle_group FROM exersise 
WHERE muscle_group = ($1)
ORDER BY exersise_name
`

func (q *Queries) GetMuscleGroupExersises(ctx context.Context, muscleGroup string) ([]Exersise, error) {
	rows, err := q.db.QueryContext(ctx, getMuscleGroupExersises, muscleGroup)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Exersise{}
	for rows.Next() {
		var i Exersise
		if err := rows.Scan(&i.ID, &i.ExersiseName, &i.MuscleGroup); err != nil {
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

const listExersises = `-- name: ListExersises :many
SELECT id, exersise_name, muscle_group FROM exersise
ORDER BY exersise_name
LIMIT $1
OFFSET $2
`

type ListExersisesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListExersises(ctx context.Context, arg ListExersisesParams) ([]Exersise, error) {
	rows, err := q.db.QueryContext(ctx, listExersises, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Exersise{}
	for rows.Next() {
		var i Exersise
		if err := rows.Scan(&i.ID, &i.ExersiseName, &i.MuscleGroup); err != nil {
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

const updateExersiseMuscleGroup = `-- name: UpdateExersiseMuscleGroup :exec
UPDATE exersise SET
muscle_group = ($1)
WHERE muscle_group = ($2)
`

type UpdateExersiseMuscleGroupParams struct {
	MuscleGroup   string `json:"muscle_group"`
	MuscleGroup_2 string `json:"muscle_group_2"`
}

func (q *Queries) UpdateExersiseMuscleGroup(ctx context.Context, arg UpdateExersiseMuscleGroupParams) error {
	_, err := q.db.ExecContext(ctx, updateExersiseMuscleGroup, arg.MuscleGroup, arg.MuscleGroup_2)
	return err
}

const updateExersiseName = `-- name: UpdateExersiseName :exec
UPDATE exersise SET
exersise_name = ($1)
WHERE exersise_name = ($2)
`

type UpdateExersiseNameParams struct {
	ExersiseName   string `json:"exersise_name"`
	ExersiseName_2 string `json:"exersise_name_2"`
}

func (q *Queries) UpdateExersiseName(ctx context.Context, arg UpdateExersiseNameParams) error {
	_, err := q.db.ExecContext(ctx, updateExersiseName, arg.ExersiseName, arg.ExersiseName_2)
	return err
}
