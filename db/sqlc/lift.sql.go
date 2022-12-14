// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: lift.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const createLift = `-- name: CreateLift :one
INSERT INTO lift (
  exercise_name,
  weight_lifted,
  reps,
  user_id,
  workout_id
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id, exercise_name, weight_lifted, reps, user_id, workout_id
`

type CreateLiftParams struct {
	ExerciseName string    `json:"exercise_name"`
	WeightLifted float32   `json:"weight_lifted"`
	Reps         int16     `json:"reps"`
	UserID       uuid.UUID `json:"user_id"`
	WorkoutID    uuid.UUID `json:"workout_id"`
}

func (q *Queries) CreateLift(ctx context.Context, arg CreateLiftParams) (Lift, error) {
	row := q.db.QueryRowContext(ctx, createLift,
		arg.ExerciseName,
		arg.WeightLifted,
		arg.Reps,
		arg.UserID,
		arg.WorkoutID,
	)
	var i Lift
	err := row.Scan(
		&i.ID,
		&i.ExerciseName,
		&i.WeightLifted,
		&i.Reps,
		&i.UserID,
		&i.WorkoutID,
	)
	return i, err
}

const createLifts = `-- name: CreateLifts :many
INSERT INTO lift (
  exercise_name,
  weight_lifted,
  reps,
  user_id,
  workout_id
) VALUES (
  UNNEST($1::VARCHAR[]),
  UNNEST($2::REAL[]),
  UNNEST($3::SMALLINT[]),
  UNNEST($4::UUID[]),
  UNNEST($5::UUID[])
)
RETURNING id, exercise_name, weight_lifted, reps, user_id, workout_id
`

type CreateLiftsParams struct {
	Exercisenames []string    `json:"exercisenames"`
	Weights       []float32   `json:"weights"`
	Reps          []int16     `json:"reps"`
	UserID        []uuid.UUID `json:"user_id"`
	WorkoutID     []uuid.UUID `json:"workout_id"`
}

func (q *Queries) CreateLifts(ctx context.Context, arg CreateLiftsParams) ([]Lift, error) {
	rows, err := q.db.QueryContext(ctx, createLifts,
		pq.Array(arg.Exercisenames),
		pq.Array(arg.Weights),
		pq.Array(arg.Reps),
		pq.Array(arg.UserID),
		pq.Array(arg.WorkoutID),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Lift{}
	for rows.Next() {
		var i Lift
		if err := rows.Scan(
			&i.ID,
			&i.ExerciseName,
			&i.WeightLifted,
			&i.Reps,
			&i.UserID,
			&i.WorkoutID,
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

const deleteLift = `-- name: DeleteLift :exec
DELETE FROM lift WHERE id = $1
`

func (q *Queries) DeleteLift(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteLift, id)
	return err
}

const getLift = `-- name: GetLift :one
SELECT id, exercise_name, weight_lifted, reps, user_id, workout_id FROM lift
WHERE user_id = $1
AND id = $2
LIMIT 1
`

type GetLiftParams struct {
	UserID uuid.UUID `json:"user_id"`
	ID     uuid.UUID `json:"id"`
}

func (q *Queries) GetLift(ctx context.Context, arg GetLiftParams) (Lift, error) {
	row := q.db.QueryRowContext(ctx, getLift, arg.UserID, arg.ID)
	var i Lift
	err := row.Scan(
		&i.ID,
		&i.ExerciseName,
		&i.WeightLifted,
		&i.Reps,
		&i.UserID,
		&i.WorkoutID,
	)
	return i, err
}

const listLifts = `-- name: ListLifts :many
SELECT id, exercise_name, weight_lifted, reps, user_id, workout_id FROM lift
WHERE user_id = $1
ORDER BY exercise_name 
LIMIT $2
OFFSET $3
`

type ListLiftsParams struct {
	UserID uuid.UUID `json:"user_id"`
	Limit  int32     `json:"limit"`
	Offset int32     `json:"offset"`
}

func (q *Queries) ListLifts(ctx context.Context, arg ListLiftsParams) ([]Lift, error) {
	rows, err := q.db.QueryContext(ctx, listLifts, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Lift{}
	for rows.Next() {
		var i Lift
		if err := rows.Scan(
			&i.ID,
			&i.ExerciseName,
			&i.WeightLifted,
			&i.Reps,
			&i.UserID,
			&i.WorkoutID,
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

const listPRs = `-- name: ListPRs :many
SELECT id, exercise_name, weight_lifted, reps, user_id, workout_id FROM lift
WHERE user_id = $1
ORDER BY
  CASE
    WHEN $2 = 'weight' THEN weight_lifted
    WHEN $3 = 'reps' THEN reps
END DESC
LIMIT $4
OFFSET $5
`

type ListPRsParams struct {
	UserID  uuid.UUID   `json:"user_id"`
	Column2 interface{} `json:"column_2"`
	Column3 interface{} `json:"column_3"`
	Limit   int32       `json:"limit"`
	Offset  int32       `json:"offset"`
}

func (q *Queries) ListPRs(ctx context.Context, arg ListPRsParams) ([]Lift, error) {
	rows, err := q.db.QueryContext(ctx, listPRs,
		arg.UserID,
		arg.Column2,
		arg.Column3,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Lift{}
	for rows.Next() {
		var i Lift
		if err := rows.Scan(
			&i.ID,
			&i.ExerciseName,
			&i.WeightLifted,
			&i.Reps,
			&i.UserID,
			&i.WorkoutID,
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

const listPRsByExercise = `-- name: ListPRsByExercise :many
SELECT id, exercise_name, weight_lifted, reps, user_id, workout_id FROM lift
WHERE user_id = $1 AND exercise_name = $2
ORDER BY
  CASE
    WHEN $3 = 'weight' THEN weight_lifted
    WHEN $4 = 'reps' THEN reps 
END DESC
LIMIT $5
OFFSET $6
`

type ListPRsByExerciseParams struct {
	UserID       uuid.UUID   `json:"user_id"`
	ExerciseName string      `json:"exercise_name"`
	Column3      interface{} `json:"column_3"`
	Column4      interface{} `json:"column_4"`
	Limit        int32       `json:"limit"`
	Offset       int32       `json:"offset"`
}

func (q *Queries) ListPRsByExercise(ctx context.Context, arg ListPRsByExerciseParams) ([]Lift, error) {
	rows, err := q.db.QueryContext(ctx, listPRsByExercise,
		arg.UserID,
		arg.ExerciseName,
		arg.Column3,
		arg.Column4,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Lift{}
	for rows.Next() {
		var i Lift
		if err := rows.Scan(
			&i.ID,
			&i.ExerciseName,
			&i.WeightLifted,
			&i.Reps,
			&i.UserID,
			&i.WorkoutID,
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

const listPRsByMuscleGroup = `-- name: ListPRsByMuscleGroup :many
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
OFFSET $6
`

type ListPRsByMuscleGroupParams struct {
	MuscleGroup string      `json:"muscle_group"`
	UserID      uuid.UUID   `json:"user_id"`
	Column3     interface{} `json:"column_3"`
	Column4     interface{} `json:"column_4"`
	Limit       int32       `json:"limit"`
	Offset      int32       `json:"offset"`
}

type ListPRsByMuscleGroupRow struct {
	ID           uuid.UUID `json:"id"`
	ExerciseName string    `json:"exercise_name"`
	WeightLifted float32   `json:"weight_lifted"`
	Reps         int16     `json:"reps"`
	MuscleGroup  string    `json:"muscle_group"`
}

func (q *Queries) ListPRsByMuscleGroup(ctx context.Context, arg ListPRsByMuscleGroupParams) ([]ListPRsByMuscleGroupRow, error) {
	rows, err := q.db.QueryContext(ctx, listPRsByMuscleGroup,
		arg.MuscleGroup,
		arg.UserID,
		arg.Column3,
		arg.Column4,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListPRsByMuscleGroupRow{}
	for rows.Next() {
		var i ListPRsByMuscleGroupRow
		if err := rows.Scan(
			&i.ID,
			&i.ExerciseName,
			&i.WeightLifted,
			&i.Reps,
			&i.MuscleGroup,
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

const updateLift = `-- name: UpdateLift :one
UPDATE lift SET
weight_lifted = COALESCE(NULLIF($1, 0::REAL), weight_lifted),
reps = COALESCE(NULLIF($2, 0::SMALLINT), reps)
WHERE id = $3
RETURNING id, exercise_name, weight_lifted, reps, user_id, workout_id
`

type UpdateLiftParams struct {
	Column1 interface{} `json:"column_1"`
	Column2 interface{} `json:"column_2"`
	ID      uuid.UUID   `json:"id"`
}

func (q *Queries) UpdateLift(ctx context.Context, arg UpdateLiftParams) (Lift, error) {
	row := q.db.QueryRowContext(ctx, updateLift, arg.Column1, arg.Column2, arg.ID)
	var i Lift
	err := row.Scan(
		&i.ID,
		&i.ExerciseName,
		&i.WeightLifted,
		&i.Reps,
		&i.UserID,
		&i.WorkoutID,
	)
	return i, err
}
