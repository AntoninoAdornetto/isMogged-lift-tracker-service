package db

import (
	"context"
	"strconv"
	"testing"

	"github.com/AntoninoAdornetto/isMogged-lift-tracker-service/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func GenerateRandLift(t *testing.T) Lift {
	exercise := GenerateRandomExercise(t)
	workout := GenerateRandWorkout(t)

	lift, err := testQueries.CreateLift(context.Background(), CreateLiftParams{
		ExerciseName: exercise.Name,
		WeightLifted: float32(util.RandomInt(100, 250)),
		Reps:         int16(util.RandomInt(4, 12)),
		UserID:       workout.UserID,
		WorkoutID:    workout.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, lift)
	require.NotNil(t, lift.ID)
	require.NotNil(t, lift.WorkoutID)
	require.NotNil(t, lift.UserID)
	require.NotNil(t, lift.Reps)
	require.NotNil(t, lift.WeightLifted)
	require.NotNil(t, lift.ExerciseName)
	return lift
}

func TestCreateLift(t *testing.T) {
	GenerateRandLift(t)
}

func TestCreateLifts(t *testing.T) {
	n := 5

	allLifts := CreateLiftsParams{
		Exercisenames: make([]string, n),
		UserID:        make([]uuid.UUID, n),
		WorkoutID:     make([]uuid.UUID, n),
		Reps:          make([]int16, n),
		Weights:       make([]float32, n),
	}

	workout := GenerateRandWorkout(t)

	for i := 0; i < n; i++ {
		ex := GenerateRandomExercise(t)

		allLifts.UserID[i] = workout.UserID
		allLifts.WorkoutID[i] = workout.ID
		allLifts.Reps[i] = int16(util.RandomInt(6, 12))
		allLifts.Weights[i] = float32(util.RandomInt(100, 220))
		allLifts.Exercisenames[i] = ex.Name
	}

	lifts, err := testQueries.CreateLifts(context.Background(), allLifts)
	require.NoError(t, err)
	require.Len(t, lifts, n)

	for i, v := range lifts {
		require.Equal(t, allLifts.UserID[i], workout.UserID)
		require.Equal(t, allLifts.WorkoutID[i], workout.ID)
		require.Equal(t, allLifts.Exercisenames[i], v.ExerciseName)
		require.Equal(t, allLifts.Weights[i], v.WeightLifted)
		require.Equal(t, allLifts.Reps[i], v.Reps)
	}
}

func TestGetLift(t *testing.T) {
	lift := GenerateRandLift(t)

	query, err := testQueries.GetLift(context.Background(), GetLiftParams{
		UserID: lift.UserID,
		ID:     lift.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, query)
	require.Equal(t, lift.ID, query.ID)
	require.Equal(t, lift.ExerciseName, query.ExerciseName)
	require.Equal(t, lift.WeightLifted, query.WeightLifted)
	require.Equal(t, lift.Reps, query.Reps)
	require.Equal(t, lift.UserID, query.UserID)
	require.Equal(t, lift.WorkoutID, query.WorkoutID)
}

func TestListLifts(t *testing.T) {
	exercise := GenerateRandomExercise(t)
	workout := GenerateRandWorkout(t)
	n := 5
	for i := 0; i < n; i++ {
		_, err := testQueries.CreateLift(context.Background(), CreateLiftParams{
			UserID:       workout.UserID,
			WorkoutID:    workout.ID,
			Reps:         int16(5 + i),
			WeightLifted: float32(150 + i),
			ExerciseName: exercise.Name,
		})
		require.NoError(t, err)
	}

	query, err := testQueries.ListLifts(context.Background(), ListLiftsParams{
		Offset: 0,
		Limit:  int32(n),
		UserID: workout.UserID,
	})
	require.NoError(t, err)
	require.Len(t, query, n)

	for _, v := range query {
		require.NotNil(t, v.UserID)
		require.NotNil(t, v.ExerciseName)
		require.NotNil(t, v.WeightLifted)
		require.NotNil(t, v.Reps)
		require.NotNil(t, v.WorkoutID)
		_ = testQueries.DeleteLift(context.Background(), v.ID)
	}
}

func TestListPRs(t *testing.T) {
	exercise := GenerateRandomExercise(t)
	workout := GenerateRandWorkout(t)

	n := 5
	reps := make([]int16, n)
	weight := make([]float32, n)

	for i := 0; i < n; i++ {
		reps[i] = int16(12 - i)
		weight[i] = float32(200 - i)
		_, err := testQueries.CreateLift(context.Background(), CreateLiftParams{
			ExerciseName: exercise.Name,
			WorkoutID:    workout.ID,
			UserID:       workout.UserID,
			WeightLifted: float32(weight[i]),
			Reps:         int16(reps[i]),
		})
		require.NoError(t, err)
	}

	orderByWeight, err := testQueries.ListPRs(context.Background(), ListPRsParams{
		UserID:  workout.UserID,
		Column3: "weight",
		Limit:   int32(n),
		Offset:  0,
	})
	require.NoError(t, err)
	require.Len(t, orderByWeight, n)

	for i, v := range orderByWeight {
		require.Equal(t, weight[i], v.WeightLifted)
		require.Equal(t, exercise.Name, v.ExerciseName)
		require.Equal(t, workout.ID, v.WorkoutID)
	}

	orderByReps, err := testQueries.ListPRs(context.Background(), ListPRsParams{
		UserID:  workout.UserID,
		Column3: "weight",
		Limit:   int32(n),
		Offset:  0,
	})
	require.NoError(t, err)
	require.Len(t, orderByReps, n)

	for i, v := range orderByReps {
		require.Equal(t, reps[i], v.Reps)
	}

	defer func() {
		_ = testQueries.DeleteWorkout(context.Background(), workout.ID)
	}()
}

func TestListPRsByExercise(t *testing.T) {
	exercise := GenerateRandomExercise(t)
	workout := GenerateRandWorkout(t)

	n := 5
	reps := make([]int16, n)
	weight := make([]float32, n)

	for i := 0; i < n; i++ {
		reps[i] = int16(12 - i)
		weight[i] = float32(200 - i)
		_, err := testQueries.CreateLift(context.Background(), CreateLiftParams{
			ExerciseName: exercise.Name,
			WorkoutID:    workout.ID,
			UserID:       workout.UserID,
			WeightLifted: float32(weight[i]),
			Reps:         int16(reps[i]),
		})
		require.NoError(t, err)
	}

	orderByWeight, err := testQueries.ListPRsByExercise(context.Background(), ListPRsByExerciseParams{
		ExerciseName: exercise.Name,
		UserID:       workout.UserID,
		Column3:      "weight",
		Limit:        int32(n),
		Offset:       0,
	})
	require.NoError(t, err)
	require.Len(t, orderByWeight, n)

	for i, v := range orderByWeight {
		require.Equal(t, weight[i], v.WeightLifted)
		require.Equal(t, exercise.Name, v.ExerciseName)
		require.Equal(t, workout.ID, v.WorkoutID)
	}

	orderByReps, err := testQueries.ListPRsByExercise(context.Background(), ListPRsByExerciseParams{
		ExerciseName: exercise.Name,
		UserID:       workout.UserID,
		Column3:      "weight",
		Limit:        int32(n),
		Offset:       0,
	})
	require.NoError(t, err)
	require.Len(t, orderByReps, n)

	for i, v := range orderByReps {
		require.Equal(t, reps[i], v.Reps)
	}

	defer func() {
		_ = testQueries.DeleteWorkout(context.Background(), workout.ID)
	}()
}

func TestListPRsByMuscleGroup(t *testing.T) {
	exercise := GenerateRandomExercise(t)
	exercise_2 := GenerateRandomExercise(t)
	workout := GenerateRandWorkout(t)

	n := 5
	reps := make([]int16, n)
	weight := make([]float32, n)

	for i := 0; i < n; i++ {
		reps[i] = int16(12 - i)
		weight[i] = float32(200 - i)
		_, err := testQueries.CreateLift(context.Background(), CreateLiftParams{
			ExerciseName: exercise.Name,
			WorkoutID:    workout.ID,
			UserID:       workout.UserID,
			WeightLifted: float32(weight[i]),
			Reps:         int16(reps[i]),
		})
		require.NoError(t, err)

		_, err = testQueries.CreateLift(context.Background(), CreateLiftParams{
			ExerciseName: exercise_2.Name,
			WorkoutID:    workout.ID,
			UserID:       workout.UserID,
			WeightLifted: float32(weight[i]),
			Reps:         int16(reps[i]),
		})
		require.NoError(t, err)
	}

	orderByWeightGroup1, err := testQueries.ListPRsByMuscleGroup(context.Background(), ListPRsByMuscleGroupParams{
		MuscleGroup: exercise.MuscleGroup,
		UserID:      workout.UserID,
		Column3:     "weight",
		Offset:      0,
		Limit:       int32(n),
	})
	require.NoError(t, err)

	for i, v := range orderByWeightGroup1 {
		require.Equal(t, weight[i], v.WeightLifted)
		require.Equal(t, exercise.MuscleGroup, v.MuscleGroup)
	}

	orderByWeightGroup2, err := testQueries.ListPRsByMuscleGroup(context.Background(), ListPRsByMuscleGroupParams{
		MuscleGroup: exercise_2.MuscleGroup,
		UserID:      workout.UserID,
		Column3:     "weight",
		Offset:      0,
		Limit:       int32(n),
	})
	require.NoError(t, err)

	for i, v := range orderByWeightGroup2 {
		require.Equal(t, weight[i], v.WeightLifted)
		require.Equal(t, exercise_2.MuscleGroup, v.MuscleGroup)
	}

	orderByRepsGroup1, err := testQueries.ListPRsByMuscleGroup(context.Background(), ListPRsByMuscleGroupParams{
		MuscleGroup: exercise.MuscleGroup,
		UserID:      workout.UserID,
		Column3:     "reps",
		Offset:      0,
		Limit:       int32(n),
	})

	for i, v := range orderByRepsGroup1 {
		require.Equal(t, reps[i], v.Reps)
		require.Equal(t, exercise.MuscleGroup, v.MuscleGroup)
	}

	orderByRepsGroup2, err := testQueries.ListPRsByMuscleGroup(context.Background(), ListPRsByMuscleGroupParams{
		MuscleGroup: exercise_2.MuscleGroup,
		UserID:      workout.UserID,
		Column3:     "reps",
		Offset:      0,
		Limit:       int32(n),
	})

	for i, v := range orderByRepsGroup2 {
		require.Equal(t, reps[i], v.Reps)
		require.Equal(t, exercise_2.MuscleGroup, v.MuscleGroup)
	}
}

func TestUpdateLift(t *testing.T) {
	patchedWeightStr := "20.5"
	patchWeightVal, _ := strconv.ParseFloat(patchedWeightStr, 32)
	patchedRepsStr := "50"
	patchedRepsVal, _ := strconv.Atoi(patchedRepsStr)

	weightLift := GenerateRandLift(t)

	patchWeightRes, err := testQueries.UpdateLift(context.Background(), UpdateLiftParams{
		ID:      weightLift.ID,
		Column1: patchedWeightStr,
		Column2: 0,
	})
	require.NoError(t, err)
	require.Equal(t, float32(patchWeightVal), patchWeightRes.WeightLifted)
	require.Equal(t, weightLift.Reps, patchWeightRes.Reps)

	repLift := GenerateRandLift(t)

	patchedRepRes, err := testQueries.UpdateLift(context.Background(), UpdateLiftParams{
		ID:      repLift.ID,
		Column1: 0,
		Column2: patchedRepsStr,
	})
	require.NoError(t, err)
	require.Equal(t, int16(patchedRepsVal), patchedRepRes.Reps)
	require.Equal(t, repLift.WeightLifted, patchedRepRes.WeightLifted)
}

func TestDeleteLift(t *testing.T) {
	lift := GenerateRandLift(t)

	err := testQueries.DeleteLift(context.Background(), lift.ID)
	require.NoError(t, err)

	_, err = testQueries.GetLift(context.Background(), GetLiftParams{
		ID:     lift.ID,
		UserID: lift.UserID,
	})
	require.Error(t, err)
}
