package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateRandomSet(t *testing.T) Set {
	acc := GenerateRandAccount(t)
	s, err := testQueries.CreateSet(context.Background(), acc.ID)
	require.NoError(t, err)
	require.NotNil(t, s.ID)
	require.NotNil(t, s.UserID)
	return s
}

func TestCreateSet(t *testing.T) {
	s := CreateRandomSet(t)
	testQueries.DeleteSet(context.Background(), s.ID)
}

func TestDeleteSet(t *testing.T) {
	s := CreateRandomSet(t)
	testQueries.DeleteSet(context.Background(), s.ID)
	query, err := testQueries.GetSet(context.Background(), s.ID)
	require.Error(t, err)
	require.Empty(t, query)
}

func TestGetLiftSets(t *testing.T) {
	var args []Lift

	acc := GenerateRandAccount(t)
	mg := CreateRandMuscleGroup(t, "Chest-Testing-Sets")

	exersiseArgs := CreateExerciseParams{
		ExerciseName: "Chest-Press-Test",
		MuscleGroup:  mg.GroupName,
	}

	ex, err := testQueries.CreateExercise(context.Background(), exersiseArgs)
	require.NoError(t, err)
	require.NotEmpty(t, ex)

	s := CreateRandomSet(t)

	for i := 0; i < 3; i++ {
		liftArgs := CreateLiftParams{
			ExerciseName: ex.ExerciseName,
			Weight:       float32(100 + i),
			Reps:         int32(i + 1),
			UserID:       acc.ID,
			SetID:        s.ID,
		}

		l, err := testQueries.CreateLift(context.Background(), liftArgs)
		require.NoError(t, err)
		require.NotEmpty(t, l.Reps)
		require.NotEmpty(t, l.Weight)
		require.NotEmpty(t, l.ID)
		require.NotEmpty(t, l.SetID)
		args = append(args, l)
	}

	liftSets, err := testQueries.GetLiftSets(context.Background(), s.ID)
	require.NoError(t, err)

	for _, set := range liftSets {
		require.NotEmpty(t, set.ExerciseName)
		require.NotEmpty(t, set.Weight)
		require.NotEmpty(t, set.Reps)
		require.NotEmpty(t, set.DateLifted)
		require.NotEmpty(t, set.SetID)
	}

	require.GreaterOrEqual(t, len(args), 3)
	testQueries.DeleteGroup(context.Background(), mg.GroupName)
	testQueries.DeleteAccount(context.Background(), acc.ID)
	testQueries.DeleteSet(context.Background(), s.ID)
}
