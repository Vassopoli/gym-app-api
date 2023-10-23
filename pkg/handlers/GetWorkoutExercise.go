package handlers

import (
	"fmt"

	"vassopoli.com/gym-api/pkg/models"
)

func (h handler) GetWorkoutExercise(idWorkout int, idExercise int) models.WorkoutExercise {
	var workoutExercise models.WorkoutExercise

	if result := h.DB.
		Where("id_workout = ? AND id_exercise = ?", idWorkout, idExercise).
		First(&workoutExercise); result.Error != nil {
		fmt.Println(result.Error)
	}

	fmt.Println(workoutExercise)
	return workoutExercise
}
