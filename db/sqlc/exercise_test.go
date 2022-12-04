package db

import (
	"context"
	"testing"

	"github.com/AntoninoAdornetto/lift_tracker/util"
	"github.com/stretchr/testify/require"
)

func GenerateRandomExercise(t *testing.T) Exercise {
	muscleGroup := GenerateRandMuscleGroup(t)
	category := GenerateRandomCategory(t)
	exerciseName := util.RandomString(5)

	exercise, err := testQueries.CreateExercise(context.Background(), CreateExerciseParams{
		Name:        exerciseName,
		MuscleGroup: muscleGroup.Name,
		Category:    category.Name,
	})
	require.NoError(t, err)
	require.NotEmpty(t, exercise)
	require.NotNil(t, exercise.Name)
	require.NotNil(t, exercise.Category)
	require.NotNil(t, exercise.MuscleGroup)
	require.NotNil(t, exercise.ID)

	return exercise
}

func TestCreateExercise(t *testing.T) {
	GenerateRandomExercise(t)
}

func TestGetExercise(t *testing.T) {
	exercise := GenerateRandomExercise(t)

	query, err := testQueries.GetExercise(context.Background(), exercise.Name)
	require.NoError(t, err)
	require.NotEmpty(t, query)
	require.Equal(t, exercise.Name, query.Name)
}

func TestListExercises(t *testing.T) {
	n := 5
	exercises := make([]Exercise, n)
	for i := n; i < n; i++ {
		exercises[i] = GenerateRandomExercise(t)
	}

	query, err := testQueries.ListExercises(context.Background(), ListExercisesParams{
		Limit:  int32(n),
		Offset: 0,
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(query), n)
	for i := 0; i < n; i++ {
		require.NotEmpty(t, query[i])
	}

	for _, v := range exercises {
		_ = testQueries.DeleteExercise(context.Background(), v.Name)
	}
}

func TestListByMuscleGroup(t *testing.T) {
	exercise := GenerateRandomExercise(t)

	query, err := testQueries.ListByMuscleGroup(context.Background(), ListByMuscleGroupParams{
		Limit:       5,
		Offset:      0,
		MuscleGroup: exercise.MuscleGroup,
	})
	require.NoError(t, err)
	for _, v := range query {
		require.Equal(t, exercise.MuscleGroup, v.MuscleGroup)
		require.NotNil(t, v.Name)
	}
}

func TestUpdateExerciseName(t *testing.T) {
	exercise := GenerateRandomExercise(t)
	newName := util.RandomString(10)

	err := testQueries.UpdateExerciseName(context.Background(), UpdateExerciseNameParams{
		Name:   newName,
		Name_2: exercise.Name,
	})
	require.NoError(t, err)

	query, err := testQueries.GetExercise(context.Background(), newName)
	require.NoError(t, err)
	require.NotEmpty(t, query)
	require.Equal(t, newName, query.Name)
}

func TestUpdateExerciseMuscleGroup(t *testing.T) {
	patchMuscleGroup := GenerateRandMuscleGroup(t)
	exercise := GenerateRandomExercise(t)

	err := testQueries.UpdateMuscleGroup(context.Background(), UpdateMuscleGroupParams{
		MuscleGroup: patchMuscleGroup.Name,
		Name:        exercise.Name,
	})
	require.NoError(t, err)

	query, err := testQueries.GetExercise(context.Background(), exercise.Name)
	require.NoError(t, err)
	require.NotEmpty(t, query)
	require.Equal(t, patchMuscleGroup.Name, query.MuscleGroup)
}

func TestDeleteExercise(t *testing.T) {
	exercise := GenerateRandomExercise(t)

	err := testQueries.DeleteExercise(context.Background(), exercise.Name)
	require.NoError(t, err)

	query, err := testQueries.GetExercise(context.Background(), exercise.Name)
	require.Error(t, err)
	require.Empty(t, query)
}
