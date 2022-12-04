package db

import (
	"context"
	"testing"

	"github.com/AntoninoAdornetto/lift_tracker/util"
	"github.com/stretchr/testify/require"
)

func GenerateRandMuscleGroup(t *testing.T) MuscleGroup {
	var muscleGroup MuscleGroup
	muscleGroupName := util.RandomString(5)

	muscleGroup, err := testQueries.CreateMuscleGroup(context.Background(), muscleGroupName)
	require.NoError(t, err)
	require.NotNil(t, muscleGroup.ID)
	require.NotNil(t, muscleGroup.Name)
	return muscleGroup
}

func TestCreateMuscleGroup(t *testing.T) {
	GenerateRandMuscleGroup(t)
}

func TestGetMuscleGroup(t *testing.T) {
	muscleGroup := GenerateRandMuscleGroup(t)

	query, err := testQueries.GetMuscleGroup(context.Background(), muscleGroup.Name)
	require.NoError(t, err)
	require.NotEmpty(t, query)
}

func TestListMuscleGroups(t *testing.T) {
	n := 5
	muscleGroups := make([]MuscleGroup, n)
	for i := 0; i < n; i++ {
		muscleGroups[i] = GenerateRandMuscleGroup(t)
	}

	query, err := testQueries.GetMuscleGroups(context.Background())
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(query), n)

	for i := 0; i < n; i++ {
		_, _ = testQueries.DeleteGroup(context.Background(), muscleGroups[i].Name)
	}
}

func TestUpdateMuscleGroup(t *testing.T) {
	muscleGroup := GenerateRandMuscleGroup(t)
	newName := util.RandomString(5)

	patch, err := testQueries.UpdateGroup(context.Background(), UpdateGroupParams{
		Name:   newName,
		Name_2: muscleGroup.Name,
	})
	require.NoError(t, err)
	require.Equal(t, newName, patch.Name)
}

func TestDeleteMuscleGroup(t *testing.T) {
	muscleGroup := GenerateRandMuscleGroup(t)

	d, err := testQueries.DeleteGroup(context.Background(), muscleGroup.Name)
	require.NoError(t, err)
	require.NotEmpty(t, d)
}
