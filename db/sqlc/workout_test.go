package db

import (
	"context"
	"testing"
	"time"

	"github.com/AntoninoAdornetto/lift_tracker/util"
	"github.com/stretchr/testify/require"
)

func GenerateRandWorkout(t *testing.T) Workout {
	account := GenerateRandAccount(t)
	startTime := time.Now().UnixMilli()

	workout, err := testQueries.CreateWorkout(context.Background(), CreateWorkoutParams{
		UserID:    account.ID,
		StartTime: util.FormatMSEpoch(startTime),
	})
	require.NoError(t, err)
	require.NotEmpty(t, workout)
	require.NotNil(t, workout.ID)
	require.NotNil(t, workout.StartTime)
	require.NotNil(t, workout.FinishTime)
	require.NotNil(t, workout.UserID)
	return workout
}

func TestCreateWorkout(t *testing.T) {
	GenerateRandWorkout(t)
}

func TestGetWorkout(t *testing.T) {
	lift := GenerateRandLift(t)

	workout, err := testQueries.GetWorkout(context.Background(), GetWorkoutParams{
		ID:     lift.WorkoutID,
		UserID: lift.UserID,
	})
	require.NoError(t, err)
	for _, v := range workout {
		require.Equal(t, lift.UserID, v.UserID)
		require.Equal(t, lift.WorkoutID, v.ID)
		require.NotNil(t, v.FinishTime)
		require.NotNil(t, v.StartTime)
		require.NotNil(t, v.Reps)
		require.NotNil(t, v.WeightLifted)
	}
}

func TestUpdateWorkoutEndTime(t *testing.T) {
	workout := GenerateRandWorkout(t)
	hour := int64(60*60*1000) + workout.StartTime.UnixMilli()
	end := util.FormatMSEpoch(hour)
	date1 := time.Date(end.Year(), end.Month(), end.Day(), end.Hour(), end.Minute(), end.Second(), end.Nanosecond(), time.UTC)
	require.Greater(t, end.UnixMilli(), workout.StartTime.UnixMilli())

	patched, err := testQueries.UpdateFinishTime(context.Background(), UpdateFinishTimeParams{
		ID:         workout.ID,
		FinishTime: end,
	})
	date2 := time.Date(patched.FinishTime.Year(), patched.FinishTime.Month(), patched.FinishTime.Day(), patched.FinishTime.Hour(), patched.FinishTime.Minute(), patched.FinishTime.Second(), patched.FinishTime.Nanosecond(), time.UTC)
	require.NoError(t, err)
	require.NotEmpty(t, patched)
	require.Equal(t, true, date1.Equal(date2))
}
