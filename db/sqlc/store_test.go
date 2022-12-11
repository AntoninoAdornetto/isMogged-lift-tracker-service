package db

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"
)

var muscleGroups = [7]string{"Chest", "Back", "Shoulders", "Legs", "Calves", "Biceps", "Triceps"}
var categories = [5]string{"Barbell", "Dumbbell", "Cable", "Machine", "Calisthenics"}

type JSONExerciseInfo struct {
	Name        string
	MuscleGroup string
	Category    string
}

type JSONCreateExercise struct {
	Exercises []JSONExerciseInfo
}

func createRealExercises() JSONCreateExercise {
	root := "../../"
	filePath := filepath.Join(root, "exercises.json")

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var exercises JSONCreateExercise
	err = json.Unmarshal(data, &exercises)
	if err != nil {
		panic(err)
	}

	return exercises
}

type JSONLift struct {
	ExerciseName string
	Reps         int16
	WeightLifted float32
}

type JSONCreateLift struct {
	Lifts []JSONLift
}

func createRealLifts() JSONCreateLift {
	root := "../../"
	filePath := filepath.Join(root, "lifts.json")

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var lifts JSONCreateLift
	err = json.Unmarshal(data, &lifts)
	if err != nil {
		panic(err)
	}

	return lifts
}

func setData(t *testing.T) {
	exercises := createRealExercises().Exercises
	for _, v := range muscleGroups {
		_, _ = testQueries.CreateMuscleGroup(context.Background(), v)
	}

	for _, v := range categories {
		_, _ = testQueries.CreateCategory(context.Background(), v)
	}

	for _, v := range exercises {
		_, _ = testQueries.CreateExercise(context.Background(), CreateExerciseParams{
			Name:        v.Name,
			Category:    v.Category,
			MuscleGroup: v.MuscleGroup,
		})
	}

	defer func() {
		workout := GenerateRandWorkout(t)
		lifts := createRealLifts().Lifts

		for _, v := range lifts {
			testQueries.CreateLift(context.Background(), CreateLiftParams{
				WorkoutID:    workout.ID,
				ExerciseName: v.ExerciseName,
				Reps:         v.Reps,
				WeightLifted: v.WeightLifted,
				UserID:       workout.UserID,
			})
		}
	}()
}

func TestJSONWorkoutData(t *testing.T) {
	setData(t)
}
