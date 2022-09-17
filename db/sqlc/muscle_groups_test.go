package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandMuscleGroup(t *testing.T, n string) MuscleGroup {
	args := MuscleGroup{
		GroupName: n,
	}

	entry, err := testQueries.CreateMuscleGroup(context.Background(), args.GroupName)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.NotNil(t, entry.ID)
	require.Equal(t, args.GroupName, entry.GroupName)

	return entry
}

func TestCreateMuscleGroup(t *testing.T) {
	n := "Chest"
	createRandMuscleGroup(t, "Chest")
	testQueries.DeleteGroup(context.Background(), n)
}

func TestDeleteMuscleGroup(t *testing.T) {
	n := "Back"
	mg := createRandMuscleGroup(t, n)
	require.NotNil(t, mg.GroupName)

	testQueries.DeleteGroup(context.Background(), n)

	query, err := testQueries.GetMuscleGroup(context.Background(), n)
	require.Error(t, err)
	require.Empty(t, query.GroupName)
}

func TestGetMuscleGroup(t *testing.T) {
	n := "Shoulders"
	mg := createRandMuscleGroup(t, n)
	require.NotNil(t, mg.GroupName)

	query, err := testQueries.GetMuscleGroup(context.Background(), n)
	require.NoError(t, err)
	require.Equal(t, n, query.GroupName)

	testQueries.DeleteGroup(context.Background(), n)
}

func TestGetMuscleGroups(t *testing.T) {
	c := "Chest"
	s := "Shoulders"
	l := "Legs"

	createRandMuscleGroup(t, c)
	createRandMuscleGroup(t, s)
	createRandMuscleGroup(t, l)

	query, err := testQueries.GetMuscleGroups(context.Background())
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(query), 3)

	testQueries.DeleteGroup(context.Background(), c)
	testQueries.DeleteGroup(context.Background(), s)
	testQueries.DeleteGroup(context.Background(), l)
}

func TestUpdateMuscleGroup(t *testing.T) {
	c := "triceps"
	u := "Triceps"

	createRandMuscleGroup(t, c)

	query, err := testQueries.GetMuscleGroup(context.Background(), c)
	require.NoError(t, err)
	require.Equal(t, query.GroupName, c)

	args := UpdateGroupParams{
		GroupName:   u,
		GroupName_2: c,
	}

	testQueries.UpdateGroup(context.Background(), args)
	patched, err := testQueries.GetMuscleGroup(context.Background(), u)
	require.NoError(t, err)
	require.Equal(t, patched.GroupName, u)

	testQueries.DeleteGroup(context.Background(), u)
}
