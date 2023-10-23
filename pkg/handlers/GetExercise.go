package handlers

import (
	"fmt"

	"gorm.io/gorm/clause"
	"vassopoli.com/gym-api/pkg/models"
)

func (h handler) GetWorkout(date string) models.Workout {
	var workout models.Workout

	if result := h.DB.Preload(clause.Associations).
		Preload("Exercises.Muscle").
		Preload("Exercises.WorkoutExercises").
		Preload("Exercises.WorkoutExercises.AdvancedTechnique").
		Preload("Exercises.WorkoutExercises.LoadType").
		// Preload("Exercises.WorkoutExercises", "workout_exercise.id_workout = 1").
		// Joins("inner join workout_exercise on workout.id = workout_exercise.id_workout").
		// Joins("inner join exercise on exercise.id = workout_exercise.id_exercise").
		Where("date IN ?", []string{date}).
		First(&workout); result.Error != nil {
		// if result := h.DB.Preload(clause.Associations).Where("id IN (SELECT id_exercise FROM workout_exercise WHERE id_workout IN (SELECT id FROM workout WHERE date IN ?))", []string{date}).Find(&workout); result.Error != nil {
		fmt.Println(result.Error)
	}

	fmt.Println(workout)
	return workout
}
