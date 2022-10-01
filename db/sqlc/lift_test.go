package db

import (
	"context"
	"testing"

	"github.com/AntoninoAdornetto/lift_tracker/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var acc Account
var set uuid.UUID
var isSetCreated = false

func cleanUpLift() {
	testQueries.DeleteAccount(context.Background(), acc.ID)
	testQueries.DeleteSet(context.Background(), set)
	isSetCreated = false
}

func CreateRandomLift(t *testing.T) Lift {
	if acc.Lifter == "" {
		acc = GenerateRandAccount(t)
	}

	if isSetCreated == false {
		set = CreateRandomSet(t)
		isSetCreated = true
	}

	en := CreateRandomExersise(t)

	arg := CreateLiftParams{
		ExersiseName: en.ExersiseName,
		Weight:       float32(util.RandomInt(100, 200)),
		Reps:         int32(util.RandomInt(6, 12)),
		UserID:       acc.ID,
		SetID:        set,
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

func TestDeleteLift(t *testing.T) {
	l := CreateRandomLift(t)

	testQueries.DeleteLift(context.Background(), l.ID)

	query, err := testQueries.GetLift(context.Background(), l.ID)
	require.Error(t, err)
	require.Empty(t, query)
}

func TestGetRepPRs(t *testing.T) {
	newAcc := GenerateRandAccount(t)
	mg := CreateRandMuscleGroup(t, "Chesticles")
	ex, err := testQueries.CreateExersise(context.Background(), CreateExersiseParams{
		ExersiseName: "Bench Press",
		MuscleGroup:  mg.GroupName,
	})

	require.NoError(t, err)

	LIFTLEN := 5

	for i := 0; i < LIFTLEN; i++ {
		testQueries.CreateLift(context.Background(), CreateLiftParams{
			ExersiseName: ex.ExersiseName,
			Weight:       float32(util.RandomInt(100, 200)),
			Reps:         int32(i + 1),
			UserID:       newAcc.ID,
			SetID:        set,
		})
	}

	prs, err := testQueries.GetRepPRs(context.Background(), newAcc.ID)
	require.NoError(t, err)
	require.Len(t, prs, LIFTLEN)

	for i := 0; i < len(prs)-1; i++ {
		require.Greater(t, prs[i+1].Reps, prs[i].Reps)
	}

	testQueries.DeleteExersise(context.Background(), ex.ExersiseName)
	testQueries.DeleteGroup(context.Background(), mg.GroupName)
}

func TestGetWeightPRs(t *testing.T) {
	newAcc := GenerateRandAccount(t)
	mg := CreateRandMuscleGroup(t, "Chesticles")
	ex, err := testQueries.CreateExersise(context.Background(), CreateExersiseParams{
		ExersiseName: "Bench Press",
		MuscleGroup:  mg.GroupName,
	})

	require.NoError(t, err)

	LIFTLEN := 5

	for i := 0; i < LIFTLEN; i++ {
		testQueries.CreateLift(context.Background(), CreateLiftParams{
			ExersiseName: ex.ExersiseName,
			Weight:       float32(i) * 2.3,
			Reps:         int32(util.RandomInt(6, 12)),
			UserID:       newAcc.ID,
			SetID:        set,
		})
	}

	prs, err := testQueries.GetWeightPRs(context.Background(), newAcc.ID)
	require.NoError(t, err)
	require.Len(t, prs, LIFTLEN)

	for i := 0; i < len(prs)-1; i++ {
		require.Greater(t, prs[i+1].Weight, prs[i].Weight)
	}

	testQueries.DeleteExersise(context.Background(), ex.ExersiseName)
	testQueries.DeleteGroup(context.Background(), mg.GroupName)
}

func TestUpdateReps(t *testing.T) {
	l := CreateRandomLift(t)

	args := UpdateRepsParams{
		Reps:   l.Reps - 1,
		ID:     l.ID,
		UserID: l.UserID,
	}

	testQueries.UpdateReps(context.Background(), args)

	query, err := testQueries.GetLift(context.Background(), l.ID)
	require.NoError(t, err)
	require.Equal(t, l.Reps-1, query.Reps)
}

func TestUpdateWeight(t *testing.T) {
	l := CreateRandomLift(t)

	args := UpdateWeightParams{
		Weight: l.Weight - 1,
		ID:     l.ID,
		UserID: l.UserID,
	}

	testQueries.UpdateWeight(context.Background(), args)

	query, err := testQueries.GetLift(context.Background(), l.ID)
	require.NoError(t, err)
	require.Equal(t, l.Weight-1, query.Weight)
	cleanUpLift()
}
