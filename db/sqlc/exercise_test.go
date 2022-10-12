package db

import (
	"context"
	"testing"

	"github.com/AntoninoAdornetto/lift_tracker/util"
	"github.com/stretchr/testify/require"
)

const MG = "TestGroup"

func CreateRandomExercise(t *testing.T) Exercise {
	en := util.RandomString(4)

	query, err := testQueries.GetMuscleGroup(context.Background(), MG)

	if err != nil {
		CreateRandMuscleGroup(t, MG)
	}

	require.NotNil(t, query.GroupName)

	args := CreateExerciseParams{
		ExerciseName: en,
		MuscleGroup:  MG,
	}

	exersise, err := testQueries.CreateExercise(context.Background(), args)
	require.NoError(t, err)
	require.Equal(t, args.ExerciseName, exersise.ExerciseName)
	require.Equal(t, args.MuscleGroup, exersise.MuscleGroup)
	require.NotNil(t, exersise.ID)

	return exersise
}

func TestCreateExercise(t *testing.T) {
	ex := CreateRandomExercise(t)
	testQueries.DeleteExercise(context.Background(), ex.ExerciseName)
}

func TestDeleteExercise(t *testing.T) {
	ex := CreateRandomExercise(t)

	testQueries.DeleteExercise(context.Background(), ex.ExerciseName)

	query, err := testQueries.GetExercise(context.Background(), ex.ExerciseName)
	require.Error(t, err)
	require.Empty(t, query.ExerciseName)
}

func TestGetExercise(t *testing.T) {
	ex := CreateRandomExercise(t)

	query, err := testQueries.GetExercise(context.Background(), ex.ExerciseName)
	require.NoError(t, err)
	require.NotEmpty(t, query.ExerciseName)

	testQueries.DeleteExercise(context.Background(), ex.ExerciseName)
}

func TestGetExercises(t *testing.T) {
	for i := 0; i < 5; i++ {
		CreateRandomExercise(t)
	}

	args := ListExercisesParams{
		Limit:  5,
		Offset: 0,
	}

	entries, err := testQueries.ListExercises(context.Background(), args)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(entries), 5)

	for i := 0; i < len(entries); i++ {
		testQueries.DeleteExercise(context.Background(), entries[i].ExerciseName)
	}
}

func TestGetMuscleGroupExercises(t *testing.T) {
	for i := 0; i < 5; i++ {
		CreateRandomExercise(t)
	}

	entries, err := testQueries.GetMuscleGroupExercises(context.Background(), MG)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(entries), 4)

	for i := 0; i < len(entries); i++ {
		testQueries.DeleteExercise(context.Background(), entries[i].ExerciseName)
	}
}

func TestUpdateExersiseMuscleGroup(t *testing.T) {
	newMG := CreateRandMuscleGroup(t, "TestUpdateExersiseMG")
	ex := CreateRandomExercise(t)

	args := UpdateExerciseMuscleGroupParams{
		MuscleGroup:  newMG.GroupName,
		ExerciseName: ex.ExerciseName,
	}

	testQueries.UpdateExerciseMuscleGroup(context.Background(), args)

	query, err := testQueries.GetMuscleGroupExercises(context.Background(), newMG.GroupName)
	require.NoError(t, err)
	require.Len(t, query, 1)

	testQueries.DeleteExercise(context.Background(), ex.ExerciseName)
	testQueries.DeleteGroup(context.Background(), newMG.GroupName)
}

func TestUpdateExersiseName(t *testing.T) {
	newEx := "Peck Deck"
	ex := CreateRandomExercise(t)

	args := UpdateExerciseNameParams{
		ExerciseName:   newEx,
		ExerciseName_2: ex.ExerciseName,
	}

	testQueries.UpdateExerciseName(context.Background(), args)

	query, err := testQueries.GetExercise(context.Background(), newEx)
	require.NoError(t, err)
	require.Equal(t, query.ExerciseName, newEx)

	testQueries.DeleteExercise(context.Background(), newEx)
}
