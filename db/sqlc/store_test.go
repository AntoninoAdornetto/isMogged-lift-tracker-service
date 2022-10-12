package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

var repRange = int32(6)
var weightRange = float32(100)

func TestCreateNewLift(t *testing.T) {
	store := NewStore(testDB)

	n := 3
	acc := GenerateRandAccount(t)
	ex := CreateRandomExercise(t)

	errs := make(chan error)
	results := make(chan CreateNewLiftRes)

	for i := 0; i < n; i++ {

		go func() {
			res, err := store.CreateNewLift(context.Background(), CreateNewLiftReq{
				ExerciseName: ex.ExerciseName,
				Reps:         repRange,
				Weight:       weightRange,
				UserID:       acc.ID,
			})

			errs <- err
			results <- res
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		res := <-results
		require.NotEmpty(t, res)
		require.Equal(t, ex.ExerciseName, res.Lift.ExerciseName)
		require.Equal(t, repRange, res.Lift.Reps)
		require.Equal(t, weightRange, res.Lift.Weight)
		require.Equal(t, acc.ID, res.Lift.UserID)
		require.NotNil(t, res.SetId)
	}
	store.DeleteAccount(context.Background(), acc.ID)
}
