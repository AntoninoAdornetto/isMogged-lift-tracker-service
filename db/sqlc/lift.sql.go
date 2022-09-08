// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: lift.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createLift = `-- name: CreateLift :one
INSERT INTO lift (
  exersise_name,
  weight,
  reps,
  user_id
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, exersise_name, weight, reps, date_lifted, user_id
`

type CreateLiftParams struct {
	ExersiseName string        `json:"exersise_name"`
	Weight       float32       `json:"weight"`
	Reps         int32         `json:"reps"`
	UserID       uuid.NullUUID `json:"user_id"`
}

func (q *Queries) CreateLift(ctx context.Context, arg CreateLiftParams) (Lift, error) {
	row := q.db.QueryRowContext(ctx, createLift,
		arg.ExersiseName,
		arg.Weight,
		arg.Reps,
		arg.UserID,
	)
	var i Lift
	err := row.Scan(
		&i.ID,
		&i.ExersiseName,
		&i.Weight,
		&i.Reps,
		&i.DateLifted,
		&i.UserID,
	)
	return i, err
}

const deleteLift = `-- name: DeleteLift :exec
DELETE FROM lift WHERE id = $1
`

func (q *Queries) DeleteLift(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteLift, id)
	return err
}

const getHighRepLifts = `-- name: GetHighRepLifts :many
SELECT id, exersise_name, weight, reps, date_lifted, user_id FROM lift 
WHERE user_id = $1
ORDER BY reps
`

func (q *Queries) GetHighRepLifts(ctx context.Context, userID uuid.NullUUID) ([]Lift, error) {
	rows, err := q.db.QueryContext(ctx, getHighRepLifts, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Lift
	for rows.Next() {
		var i Lift
		if err := rows.Scan(
			&i.ID,
			&i.ExersiseName,
			&i.Weight,
			&i.Reps,
			&i.DateLifted,
			&i.UserID,
		); err != nil {
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

const getLargeWeightLifts = `-- name: GetLargeWeightLifts :many
SELECT id, exersise_name, weight, reps, date_lifted, user_id FROM lift 
WHERE user_id = $1
ORDER BY weight
`

func (q *Queries) GetLargeWeightLifts(ctx context.Context, userID uuid.NullUUID) ([]Lift, error) {
	rows, err := q.db.QueryContext(ctx, getLargeWeightLifts, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Lift
	for rows.Next() {
		var i Lift
		if err := rows.Scan(
			&i.ID,
			&i.ExersiseName,
			&i.Weight,
			&i.Reps,
			&i.DateLifted,
			&i.UserID,
		); err != nil {
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

const getLift = `-- name: GetLift :one
SELECT id, exersise_name, weight, reps, date_lifted, user_id FROM lift 
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetLift(ctx context.Context, id int64) (Lift, error) {
	row := q.db.QueryRowContext(ctx, getLift, id)
	var i Lift
	err := row.Scan(
		&i.ID,
		&i.ExersiseName,
		&i.Weight,
		&i.Reps,
		&i.DateLifted,
		&i.UserID,
	)
	return i, err
}

const updateReps = `-- name: UpdateReps :exec
UPDATE lift SET
reps = $1
WHERE id = $2 AND
user_id = $3
`

type UpdateRepsParams struct {
	Reps   int32         `json:"reps"`
	ID     int64         `json:"id"`
	UserID uuid.NullUUID `json:"user_id"`
}

func (q *Queries) UpdateReps(ctx context.Context, arg UpdateRepsParams) error {
	_, err := q.db.ExecContext(ctx, updateReps, arg.Reps, arg.ID, arg.UserID)
	return err
}

const updateWeight = `-- name: UpdateWeight :exec
UPDATE lift SET
weight = $1
WHERE id = $2 AND
user_id = $3
`

type UpdateWeightParams struct {
	Weight float32       `json:"weight"`
	ID     int64         `json:"id"`
	UserID uuid.NullUUID `json:"user_id"`
}

func (q *Queries) UpdateWeight(ctx context.Context, arg UpdateWeightParams) error {
	_, err := q.db.ExecContext(ctx, updateWeight, arg.Weight, arg.ID, arg.UserID)
	return err
}
