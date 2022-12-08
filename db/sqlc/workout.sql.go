// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: workout.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createWorkout = `-- name: CreateWorkout :one
INSERT INTO workout (user_id, start_time) 
VALUES ($1, $2)
RETURNING id, start_time, finish_time, user_id
`

type CreateWorkoutParams struct {
	UserID    uuid.UUID `json:"user_id"`
	StartTime time.Time `json:"start_time"`
}

func (q *Queries) CreateWorkout(ctx context.Context, arg CreateWorkoutParams) (Workout, error) {
	row := q.db.QueryRowContext(ctx, createWorkout, arg.UserID, arg.StartTime)
	var i Workout
	err := row.Scan(
		&i.ID,
		&i.StartTime,
		&i.FinishTime,
		&i.UserID,
	)
	return i, err
}

const deleteWorkout = `-- name: DeleteWorkout :exec

DELETE FROM workout
WHERE id = $1
`

// @todo
// ListWeightPRWorkouts :many
// ListRepPRWorkouts :many
func (q *Queries) DeleteWorkout(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteWorkout, id)
	return err
}

const getWorkout = `-- name: GetWorkout :many
SELECT w.id, exercise_name, weight_lifted, reps, start_time, finish_time, l.user_id
FROM workout AS w
JOIN lift AS l ON w.id = $1 AND l.user_id = $2
`

type GetWorkoutParams struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

type GetWorkoutRow struct {
	ID           uuid.UUID `json:"id"`
	ExerciseName string    `json:"exercise_name"`
	WeightLifted float32   `json:"weight_lifted"`
	Reps         int16     `json:"reps"`
	StartTime    time.Time `json:"start_time"`
	FinishTime   time.Time `json:"finish_time"`
	UserID       uuid.UUID `json:"user_id"`
}

func (q *Queries) GetWorkout(ctx context.Context, arg GetWorkoutParams) ([]GetWorkoutRow, error) {
	rows, err := q.db.QueryContext(ctx, getWorkout, arg.ID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetWorkoutRow{}
	for rows.Next() {
		var i GetWorkoutRow
		if err := rows.Scan(
			&i.ID,
			&i.ExerciseName,
			&i.WeightLifted,
			&i.Reps,
			&i.StartTime,
			&i.FinishTime,
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

const listWorkouts = `-- name: ListWorkouts :many
SELECT
w.id, exercise_name, weight_lifted, reps, start_time, finish_time, l.user_id
FROM workout AS w
JOIN lift AS l ON w.id = $1 AND l.user_id = $2
ORDER BY start_time
LIMIT $3
OFFSET $4
`

type ListWorkoutsParams struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Limit  int32     `json:"limit"`
	Offset int32     `json:"offset"`
}

type ListWorkoutsRow struct {
	ID           uuid.UUID `json:"id"`
	ExerciseName string    `json:"exercise_name"`
	WeightLifted float32   `json:"weight_lifted"`
	Reps         int16     `json:"reps"`
	StartTime    time.Time `json:"start_time"`
	FinishTime   time.Time `json:"finish_time"`
	UserID       uuid.UUID `json:"user_id"`
}

func (q *Queries) ListWorkouts(ctx context.Context, arg ListWorkoutsParams) ([]ListWorkoutsRow, error) {
	rows, err := q.db.QueryContext(ctx, listWorkouts,
		arg.ID,
		arg.UserID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListWorkoutsRow{}
	for rows.Next() {
		var i ListWorkoutsRow
		if err := rows.Scan(
			&i.ID,
			&i.ExerciseName,
			&i.WeightLifted,
			&i.Reps,
			&i.StartTime,
			&i.FinishTime,
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

const updateFinishTime = `-- name: UpdateFinishTime :one
UPDATE workout SET
finish_time = $1
WHERE id = $2
RETURNING id, start_time, finish_time, user_id
`

type UpdateFinishTimeParams struct {
	FinishTime time.Time `json:"finish_time"`
	ID         uuid.UUID `json:"id"`
}

func (q *Queries) UpdateFinishTime(ctx context.Context, arg UpdateFinishTimeParams) (Workout, error) {
	row := q.db.QueryRowContext(ctx, updateFinishTime, arg.FinishTime, arg.ID)
	var i Workout
	err := row.Scan(
		&i.ID,
		&i.StartTime,
		&i.FinishTime,
		&i.UserID,
	)
	return i, err
}
