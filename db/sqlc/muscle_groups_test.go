package db

import (
	"context"
	"testing"

	"github.com/AntoninoAdornetto/lift_tracker/util"
	"github.com/stretchr/testify/require"
)

func CreateRandMuscleGroup(t *testing.T, n string) MuscleGroup {
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
	n := util.RandomString(10)
	CreateRandMuscleGroup(t, n)
}

func TestDeleteMuscleGroup(t *testing.T) {
	n := util.RandomString(10)
	mg := CreateRandMuscleGroup(t, n)
	require.NotNil(t, mg.GroupName)

	testQueries.DeleteGroup(context.Background(), n)

	query, err := testQueries.GetMuscleGroup(context.Background(), n)
	require.Error(t, err)
	require.Empty(t, query.GroupName)
}

func TestGetMuscleGroup(t *testing.T) {
	n := util.RandomString(7)
	mg := CreateRandMuscleGroup(t, n)
	require.NotNil(t, mg.GroupName)

	query, err := testQueries.GetMuscleGroup(context.Background(), n)
	require.NoError(t, err)
	require.Equal(t, n, query.GroupName)
}

func TestGetMuscleGroups(t *testing.T) {
	c := util.RandomString(5)
	s := util.RandomString(6)
	l := util.RandomString(7)

	CreateRandMuscleGroup(t, c)
	CreateRandMuscleGroup(t, s)
	CreateRandMuscleGroup(t, l)

	query, err := testQueries.GetMuscleGroups(context.Background())
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(query), 3)
}

func TestUpdateMuscleGroup(t *testing.T) {
	c := util.RandomString(5)
	u := util.RandomString(6)

	CreateRandMuscleGroup(t, c)

	query, err := testQueries.GetMuscleGroup(context.Background(), c)
	require.NoError(t, err)
	require.Equal(t, query.GroupName, c)

	args := UpdateGroupParams{
		GroupName:   u,
		GroupName_2: c,
	}

	patch, err := testQueries.UpdateGroup(context.Background(), args)
	require.NoError(t, err)
	require.Equal(t, patch.GroupName, u)

	testQueries.DeleteGroup(context.Background(), u)
}
