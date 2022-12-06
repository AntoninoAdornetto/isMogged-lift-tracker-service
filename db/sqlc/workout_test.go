package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func GenerateRandWorkout(t *testing.T) Workout {
	account := GenerateRandAccount(t)

	workout, err := testQueries.CreateWorkout(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, workout)
	require.NotNil(t, workout.ID)
	require.NotNil(t, workout.StartTime)
	require.NotNil(t, workout.FinishTime)
	require.NotNil(t, workout.UserID)
	return workout
}

func TestCreateWorkout(t *testing.T) {
	GenerateRandWorkout(t)
}

func TestGetWorkout(t *testing.T) {
	lift := GenerateRandLift(t)

	workout, err := testQueries.GetWorkout(context.Background(), GetWorkoutParams{
		ID:     lift.WorkoutID,
		UserID: lift.UserID,
	})
	require.NoError(t, err)
	for _, v := range workout {
		require.Equal(t, lift.UserID, v.UserID)
		require.Equal(t, lift.WorkoutID, v.ID)
		require.NotNil(t, v.FinishTime)
		require.NotNil(t, v.StartTime)
		require.NotNil(t, v.Reps)
		require.NotNil(t, v.WeightLifted)
	}
}
