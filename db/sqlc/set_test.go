package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func CreateRandomSet(t *testing.T) uuid.UUID {
	s, err := testQueries.CreateSet(context.Background())
	require.NoError(t, err)
	require.NotNil(t, s.ID)
	return s
}

func TestCreateSet(t *testing.T) {
	s := CreateRandomSet(t)
	testQueries.DeleteSet(context.Background(), s)
}

func TestDeleteSet(t *testing.T) {
	s := CreateRandomSet(t)
	q := testQueries.DeleteSet(context.Background(), s)
	require.Empty(t, q) // TODO: Add better assertions.
}

func TestGetLiftSets(t *testing.T) {
	var args []Lift

	acc := GenerateRandAccount(t)
	mg := CreateRandMuscleGroup(t, "Chest-Testing-Sets")

	exersiseArgs := CreateExersiseParams{
		ExersiseName: "Chest-Press-Test",
		MuscleGroup:  mg.GroupName,
	}

	ex, err := testQueries.CreateExersise(context.Background(), exersiseArgs)
	require.NoError(t, err)
	require.NotEmpty(t, ex)

	s := CreateRandomSet(t)

	for i := 0; i < 3; i++ {
		liftArgs := CreateLiftParams{
			ExersiseName: ex.ExersiseName,
			Weight:       float32(100 + i),
			Reps:         int32(i + 1),
			UserID:       acc.ID,
			SetID:        s,
		}

		l, err := testQueries.CreateLift(context.Background(), liftArgs)
		require.NoError(t, err)
		require.NotEmpty(t, l.Reps)
		require.NotEmpty(t, l.Weight)
		require.NotEmpty(t, l.ID)
		require.NotEmpty(t, l.SetID)
		args = append(args, l)
	}

	liftSets, err := testQueries.GetLiftSets(context.Background(), s)
	require.NoError(t, err)

	for _, set := range liftSets {
		require.NotEmpty(t, set.ExersiseName)
		require.NotEmpty(t, set.Weight)
		require.NotEmpty(t, set.Reps)
		require.NotEmpty(t, set.DateLifted)
		require.NotEmpty(t, set.SetID)
	}

	require.GreaterOrEqual(t, len(args), 3)
	testQueries.DeleteGroup(context.Background(), mg.GroupName)
	testQueries.DeleteAccount(context.Background(), acc.ID)
	testQueries.DeleteSet(context.Background(), s)
}
