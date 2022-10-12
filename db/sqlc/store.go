package db

import (
	"context"
	"database/sql"
	"fmt"

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
			return fmt.Errorf("transaction error: %v\nrb err: %v\n", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type CreateNewLiftReq struct {
	ExerciseName string    `json:"exersise_name"`
	Weight       float32   `json:"weight"`
	Reps         int32     `json:"reps"`
	UserID       uuid.UUID `json:"user_id"`
}

type CreateNewLiftRes struct {
	Lift  Lift      `json:"lift"`
	SetId uuid.UUID `json:"set"`
}

func (store *Store) CreateNewLift(ctx context.Context, args CreateNewLiftReq) (CreateNewLiftRes, error) {
	var res CreateNewLiftRes

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		res.SetId, err = q.CreateSet(ctx)

		if err != nil {
			return err
		}

		res.Lift, err = q.CreateLift(ctx, CreateLiftParams{
			ExerciseName: args.ExerciseName,
			Weight:       args.Weight,
			Reps:         args.Reps,
			SetID:        res.SetId,
			UserID:       args.UserID,
		})

		if err != nil {
			return err
		}

		return nil
	})

	return res, err
}
