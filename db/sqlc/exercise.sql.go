// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: exercise.sql

package db

import (
	"context"
)

const createExercise = `-- name: CreateExercise :one
INSERT INTO exercise (
  exercise_name,
  muscle_group
) VALUES (
  $1, $2
) RETURNING id, exercise_name, muscle_group
`

type CreateExerciseParams struct {
	ExerciseName string `json:"exercise_name"`
	MuscleGroup  string `json:"muscle_group"`
}

func (q *Queries) CreateExercise(ctx context.Context, arg CreateExerciseParams) (Exercise, error) {
	row := q.db.QueryRowContext(ctx, createExercise, arg.ExerciseName, arg.MuscleGroup)
	var i Exercise
	err := row.Scan(&i.ID, &i.ExerciseName, &i.MuscleGroup)
	return i, err
}

const deleteExercise = `-- name: DeleteExercise :exec
DELETE FROM exercise WHERE exercise_name = ($1)
`

func (q *Queries) DeleteExercise(ctx context.Context, exerciseName string) error {
	_, err := q.db.ExecContext(ctx, deleteExercise, exerciseName)
	return err
}

const getExercise = `-- name: GetExercise :one
SELECT id, exercise_name, muscle_group FROM exercise
WHERE exercise_name = ($1) LIMIT 1
`

func (q *Queries) GetExercise(ctx context.Context, exerciseName string) (Exercise, error) {
	row := q.db.QueryRowContext(ctx, getExercise, exerciseName)
	var i Exercise
	err := row.Scan(&i.ID, &i.ExerciseName, &i.MuscleGroup)
	return i, err
}

const getMuscleGroupExercises = `-- name: GetMuscleGroupExercises :many
SELECT id, exercise_name, muscle_group FROM exercise 
WHERE muscle_group = ($1)
ORDER BY exercise_name
`

func (q *Queries) GetMuscleGroupExercises(ctx context.Context, muscleGroup string) ([]Exercise, error) {
	rows, err := q.db.QueryContext(ctx, getMuscleGroupExercises, muscleGroup)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Exercise{}
	for rows.Next() {
		var i Exercise
		if err := rows.Scan(&i.ID, &i.ExerciseName, &i.MuscleGroup); err != nil {
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

const listExercises = `-- name: ListExercises :many
SELECT id, exercise_name, muscle_group FROM exercise
ORDER BY exercise_name
LIMIT $1
OFFSET $2
`

type ListExercisesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListExercises(ctx context.Context, arg ListExercisesParams) ([]Exercise, error) {
	rows, err := q.db.QueryContext(ctx, listExercises, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Exercise{}
	for rows.Next() {
		var i Exercise
		if err := rows.Scan(&i.ID, &i.ExerciseName, &i.MuscleGroup); err != nil {
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

const updateExerciseMuscleGroup = `-- name: UpdateExerciseMuscleGroup :exec
UPDATE exercise SET
muscle_group = ($1)
WHERE exercise_name = ($2)
`

type UpdateExerciseMuscleGroupParams struct {
	MuscleGroup  string `json:"muscle_group"`
	ExerciseName string `json:"exercise_name"`
}

func (q *Queries) UpdateExerciseMuscleGroup(ctx context.Context, arg UpdateExerciseMuscleGroupParams) error {
	_, err := q.db.ExecContext(ctx, updateExerciseMuscleGroup, arg.MuscleGroup, arg.ExerciseName)
	return err
}

const updateExerciseName = `-- name: UpdateExerciseName :exec
UPDATE exercise SET
exercise_name = ($1)
WHERE exercise_name = ($2)
`

type UpdateExerciseNameParams struct {
	ExerciseName   string `json:"exercise_name"`
	ExerciseName_2 string `json:"exercise_name_2"`
}

func (q *Queries) UpdateExerciseName(ctx context.Context, arg UpdateExerciseNameParams) error {
	_, err := q.db.ExecContext(ctx, updateExerciseName, arg.ExerciseName, arg.ExerciseName_2)
	return err
}
