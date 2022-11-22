package db

import (
	"context"
	"testing"

	"github.com/AntoninoAdornetto/lift_tracker/util"
	"github.com/stretchr/testify/require"
)

var acc Account
var set Set
var isSetCreated = false

func CreateRandomLift(t *testing.T) Lift {
	if acc.Lifter == "" {
		acc = GenerateRandAccount(t)
	}

	if isSetCreated == false {
		set = CreateRandomSet(t)
		isSetCreated = true
	}

	en := CreateRandomExercise(t)

	arg := CreateLiftParams{
		ExerciseName: en.ExerciseName,
		Weight:       float32(util.RandomInt(100, 200)),
		Reps:         int32(util.RandomInt(6, 12)),
		UserID:       acc.ID,
		SetID:        set.ID,
	}

	l, err := testQueries.CreateLift(context.Background(), arg)
	require.NoError(t, err)
	require.NotNil(t, l.ID)
	require.NotNil(t, l.SetID)
	require.NotNil(t, l.UserID)
	require.NotNil(t, l.UserID)
	require.NotNil(t, l.Reps)
	require.NotNil(t, l.Weight)

	require.NotNil(t, acc.ID)

	return l
}

func TestCreateLift(t *testing.T) {
	CreateRandomLift(t)
}

func TestGetLift(t *testing.T) {
	l := CreateRandomLift(t)

	query, err := testQueries.GetLift(context.Background(), l.ID)
	require.NoError(t, err)
	require.NotNil(t, query.ID)
}

func TestListLifts(t *testing.T) {
	_1 := CreateRandomLift(t)
	l := 3

	for i := 0; i < l-1; i++ {
		CreateRandomLift(t)
	}

	args := ListLiftsParams{
		UserID: _1.UserID,
		Limit:  int32(l),
		Offset: 0,
	}

	query, err := testQueries.ListLifts(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, query, l)

	for _, d := range query {
		require.NotNil(t, d.ID)
		require.NotNil(t, d.UserID)
		require.NotNil(t, d.ExerciseName)
		require.NotNil(t, d.Weight)
		require.NotNil(t, d.Reps)
		require.NotNil(t, d.SetID)
		require.NotNil(t, d.DateLifted)
	}
}

func TestDeleteLift(t *testing.T) {
	l := CreateRandomLift(t)

	testQueries.DeleteLift(context.Background(), l.ID)

	query, err := testQueries.GetLift(context.Background(), l.ID)
	require.Error(t, err)
	require.Empty(t, query)
}

func TestListRepPRLifts(t *testing.T) {
	newAcc := GenerateRandAccount(t)
	mg := CreateRandMuscleGroup(t, util.RandomString(10))

	ex, err := testQueries.CreateExercise(context.Background(), CreateExerciseParams{
		ExerciseName: util.RandomString(9),
		MuscleGroup:  mg.GroupName,
	})

	require.NoError(t, err)

	LIFTLEN := 5

	for i := 0; i < LIFTLEN; i++ {
		testQueries.CreateLift(context.Background(), CreateLiftParams{
			ExerciseName: ex.ExerciseName,
			Weight:       float32(i),
			Reps:         int32(util.RandomInt(1, 100)),
			UserID:       newAcc.ID,
			SetID:        set.ID,
		})
	}

	prs, err := testQueries.ListRepPRs(context.Background(), ListRepPRsParams{
		UserID: acc.ID,
		Limit:  10,
		Offset: 0,
	})
	require.NoError(t, err)

	for i := 0; i < len(prs)-1; i++ {
		require.GreaterOrEqual(t, prs[i].Reps, prs[i+1].Reps)
		require.NotNil(t, prs[i].ExerciseName)
		require.NotNil(t, prs[i].ID)
		require.NotNil(t, prs[i].SetID)
		require.NotNil(t, prs[i].Weight)
	}

	testQueries.DeleteExercise(context.Background(), ex.ExerciseName)
	testQueries.DeleteGroup(context.Background(), mg.GroupName)
}

func TestListWeightPRLifts(t *testing.T) {
	newAcc := GenerateRandAccount(t)
	mg := CreateRandMuscleGroup(t, "Chesticles")
	ex, err := testQueries.CreateExercise(context.Background(), CreateExerciseParams{
		ExerciseName: "Bench Press",
		MuscleGroup:  mg.GroupName,
	})

	require.NoError(t, err)

	LIFTLEN := 5

	for i := 0; i < LIFTLEN; i++ {
		testQueries.CreateLift(context.Background(), CreateLiftParams{
			ExerciseName: ex.ExerciseName,
			Weight:       float32(i),
			Reps:         int32(util.RandomInt(6, 12)),
			UserID:       newAcc.ID,
			SetID:        set.ID,
		})
	}

	prs, err := testQueries.ListWeightPRLifts(context.Background(), ListWeightPRLiftsParams{
		UserID: acc.ID,
		Limit:  10,
		Offset: 0,
	})
	require.NoError(t, err)

	for i := 0; i < len(prs)-1; i++ {
		require.GreaterOrEqual(t, prs[i].Weight, prs[i+1].Weight)
	}

	testQueries.DeleteExercise(context.Background(), ex.ExerciseName)
	testQueries.DeleteGroup(context.Background(), mg.GroupName)
}

func TestUpdateReps(t *testing.T) {
	l := CreateRandomLift(t)

	args := UpdateRepsParams{
		Reps:   l.Reps - 1,
		ID:     l.ID,
		UserID: l.UserID,
	}

	patch, err := testQueries.UpdateReps(context.Background(), args)

	require.NoError(t, err)
	require.Equal(t, l.Reps-1, patch.Reps)
}

func TestUpdateLiftWeight(t *testing.T) {
	l := CreateRandomLift(t)

	args := UpdateLiftWeightParams{
		Weight: l.Weight - 1,
		ID:     l.ID,
		UserID: l.UserID,
	}

	testQueries.UpdateLiftWeight(context.Background(), args)

	query, err := testQueries.GetLift(context.Background(), l.ID)
	require.NoError(t, err)
	require.Equal(t, l.Weight-1, query.Weight)
}

func TestListNamedLiftWeightPRS(t *testing.T) {
	acc_set := CreateRandomSet(t)
	mg := CreateRandMuscleGroup(t, util.RandomString(10))
	exName := util.RandomString(9)

	ex, err := testQueries.CreateExercise(context.Background(), CreateExerciseParams{
		ExerciseName: exName,
		MuscleGroup:  mg.GroupName,
	})

	args := ListNamedLiftWeightPRsParams{
		ExerciseName: exName,
		UserID:       acc_set.UserID,
		Limit:        5,
		Offset:       0,
	}

	require.NoError(t, err)

	LIFTLEN := 5

	for i := 0; i < LIFTLEN; i++ {
		testQueries.CreateLift(context.Background(), CreateLiftParams{
			ExerciseName: ex.ExerciseName,
			Weight:       float32(i),
			Reps:         int32(util.RandomInt(1, 100)),
			UserID:       acc_set.UserID,
			SetID:        acc_set.ID,
		})
	}

	prs, err := testQueries.ListNamedLiftWeightPRs(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, prs, LIFTLEN)

	for i := 0; i < len(prs)-1; i++ {
		require.GreaterOrEqual(t, prs[i].Weight, prs[i+1].Weight)
		require.Equal(t, prs[i].ExerciseName, ex.ExerciseName)
		require.Equal(t, prs[i].SetID, acc_set.ID)
		require.Equal(t, prs[i].UserID, acc_set.UserID)
		require.NotNil(t, prs[i].Reps)
		require.NotNil(t, prs[i].DateLifted)
	}
}

func TestListMuscleGroupPRS(t *testing.T) {
	acc_set := CreateRandomSet(t)
	mg := CreateRandMuscleGroup(t, util.RandomString(10))
	exName := util.RandomString(9)

	ex, err := testQueries.CreateExercise(context.Background(), CreateExerciseParams{
		ExerciseName: exName,
		MuscleGroup:  mg.GroupName,
	})

	require.NoError(t, err)

	LIFTLEN := 5

	for i := 0; i < LIFTLEN; i++ {
		testQueries.CreateLift(context.Background(), CreateLiftParams{
			ExerciseName: ex.ExerciseName,
			Weight:       float32(i),
			Reps:         int32(util.RandomInt(1, 100)),
			UserID:       acc_set.UserID,
			SetID:        acc_set.ID,
		})
	}

	args := ListMuscleGroupPRsParams{
		UserID:      acc_set.UserID,
		MuscleGroup: mg.GroupName,
		Limit:       5,
		Offset:      0,
	}

	prs, err := testQueries.ListMuscleGroupPRs(context.Background(), args)

	require.NoError(t, err)
	require.Len(t, prs, LIFTLEN)

	for i := 0; i < LIFTLEN; i++ {
		require.Equal(t, prs[i].MuscleGroup, mg.GroupName)
	}
}
