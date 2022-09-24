package db

import (
	"context"
	"testing"

	"github.com/AntoninoAdornetto/lift_tracker/util"
	"github.com/stretchr/testify/require"
)

const MG = "TestGroup"

func createRandomExersise(t *testing.T) Exersise {
	en := util.RandomString(4)

	query, err := testQueries.GetMuscleGroup(context.Background(), MG)

	if err != nil {
		createRandMuscleGroup(t, MG)
	}

	require.NotNil(t, query.GroupName)

	args := CreateExersiseParams{
		ExersiseName: en,
		MuscleGroup:  MG,
	}

	exersise, err := testQueries.CreateExersise(context.Background(), args)
	require.NoError(t, err)
	require.Equal(t, args.ExersiseName, exersise.ExersiseName)
	require.Equal(t, args.MuscleGroup, exersise.MuscleGroup)
	require.NotNil(t, exersise.ID)

	return exersise
}

func TestCreateExersise(t *testing.T) {
	ex := createRandomExersise(t)
	testQueries.DeleteExersise(context.Background(), ex.ExersiseName)
}

func TestDeleteExersise(t *testing.T) {
	ex := createRandomExersise(t)

	testQueries.DeleteExersise(context.Background(), ex.ExersiseName)

	query, err := testQueries.GetExersise(context.Background(), ex.ExersiseName)
	require.Error(t, err)
	require.Empty(t, query.ExersiseName)
}

func TestGetExersise(t *testing.T) {
	ex := createRandomExersise(t)

	query, err := testQueries.GetExersise(context.Background(), ex.ExersiseName)
	require.NoError(t, err)
	require.NotEmpty(t, query.ExersiseName)

	testQueries.DeleteExersise(context.Background(), ex.ExersiseName)
}

func TestGetExersises(t *testing.T) {
	for i := 0; i < 5; i++ {
		createRandomExersise(t)
	}

	entries, err := testQueries.GetExersises(context.Background())
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(entries), 4)

	for i := 0; i < len(entries); i++ {
		testQueries.DeleteExersise(context.Background(), entries[i].ExersiseName)
	}
}

func TestGetMuscleGroupExersises(t *testing.T) {
	for i := 0; i < 5; i++ {
		createRandomExersise(t)
	}

	entries, err := testQueries.GetMuscleGroupExersises(context.Background(), MG)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(entries), 4)

	for i := 0; i < len(entries); i++ {
		testQueries.DeleteExersise(context.Background(), entries[i].ExersiseName)
	}
}

func TestUpdateExersiseMuscleGroup(t *testing.T) {
	newMG := createRandMuscleGroup(t, "NewGroup")
	ex := createRandomExersise(t)

	args := UpdateExersiseMuscleGroupParams{
		MuscleGroup:   newMG.GroupName,
		MuscleGroup_2: MG,
	}

	testQueries.UpdateExersiseMuscleGroup(context.Background(), args)

	query, err := testQueries.GetMuscleGroupExersises(context.Background(), newMG.GroupName)
	require.NoError(t, err)
	require.Len(t, query, 1)

	testQueries.DeleteExersise(context.Background(), ex.ExersiseName)
	testQueries.DeleteGroup(context.Background(), newMG.GroupName)
}

func TestUpdateExersiseName(t *testing.T) {
	newEx := "Peck Deck"
	ex := createRandomExersise(t)

	args := UpdateExersiseNameParams{
		ExersiseName:   newEx,
		ExersiseName_2: ex.ExersiseName,
	}

	testQueries.UpdateExersiseName(context.Background(), args)

	query, err := testQueries.GetExersise(context.Background(), newEx)
	require.NoError(t, err)
	require.Equal(t, query.ExersiseName, newEx)

	testQueries.DeleteExersise(context.Background(), newEx)
}
