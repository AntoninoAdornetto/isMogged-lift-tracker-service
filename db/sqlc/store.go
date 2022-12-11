package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/AntoninoAdornetto/lift_tracker/util"
	"github.com/google/uuid"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %v\nrb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type CreateCompleteWorkoutReq struct {
	UserId       uuid.UUID `json:"user_id"`
	StartTime    int64     `json:"start_time"`
	ExerciseName string    `json:"exercise_name"`
	Reps         int16     `json:"reps"`
	WeightLifted float32   `json:"weight_lifted"`
}

type CreateCompleteWorkoutRes = []struct {
	id            uuid.UUID
	exercise_name string
	weight_lifted float32
	reps          int16
	start_time    time.Time
	finish_time   time.Time
	userId        uuid.UUID
}

func (store *Store) CreateCompleteWorkoutTX(ctx context.Context, args CreateCompleteWorkoutReq) (CreateCompleteWorkoutRes, error) {
	var res CreateCompleteWorkoutRes
	finishTime := args.StartTime + (60 * 60 * 1000)

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		workout, err := q.CreateWorkout(ctx, CreateWorkoutParams{
			UserID:    args.UserId,
			StartTime: util.FormatMSEpoch(args.StartTime),
		})
		if err != nil {
			return err
		}

		_, err = q.CreateLift(ctx, CreateLiftParams{
			ExerciseName: args.ExerciseName,
			UserID:       args.UserId,
			WorkoutID:    workout.ID,
			Reps:         args.Reps,
			WeightLifted: args.WeightLifted,
		})

		if err != nil {
			return err
		}

		_, err = q.UpdateFinishTime(ctx, UpdateFinishTimeParams{
			ID:         workout.ID,
			FinishTime: util.FormatMSEpoch(finishTime),
		})
		if err != nil {
			return err
		}

		return nil
	})

	return res, err
}
