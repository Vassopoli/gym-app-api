package handlers

import (
	"fmt"

	"vassopoli.com/gym-api/pkg/models"
)

func (h handler) UpdateWorkoutExercise(workoutExercise models.WorkoutExercise) models.WorkoutExercise {
	if result := h.DB.
		Save(&workoutExercise); result.Error != nil {
		fmt.Println(result.Error)
	}

	fmt.Println(workoutExercise)
	return workoutExercise
}
