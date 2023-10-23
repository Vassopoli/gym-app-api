package controllers

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"vassopoli.com/gym-api/pkg/db"
	"vassopoli.com/gym-api/pkg/handlers"
	"vassopoli.com/gym-api/pkg/models"
)

var dateWorkoutCache map[string]models.Workout = make(map[string]models.Workout)

func GetWorkout(c *gin.Context) {
	dateWorkout := c.DefaultQuery("date", time.Now().Format("2006-01-02")) //TODO: Maybe move from query parameter to a path parameter, because seems like a unique resource, not a filter. Also change id on database to the date
	// dateWorkout = "2023-10-19"                                             //TODO: use as a workaround because I don't have the control on the front end yet

	//TODO: Use a real cache implementation
	currentDateWorkoutCache, ok := dateWorkoutCache[dateWorkout]
	if ok {
		c.JSON(http.StatusOK, currentDateWorkoutCache)
		return
	}

	//TODO: create just once
	DB := db.Init()
	h := handlers.New(DB)

	workout := h.GetWorkout(dateWorkout)

	//TODO: fix gambiarra - sort
	//Reordering the WorkoutExercise array, where the newest comes first, and the oldest comes last, for a matter of contract with the front end (And the most recent exercise will be the exercise of the current workout)
	for _, exercise := range workout.Exercises {
		sort.Slice(exercise.WorkoutExercises, func(i, j int) bool {
			return exercise.WorkoutExercises[i].WorkoutID > exercise.WorkoutExercises[j].WorkoutID
		})
	}

	//TODO: fix gambiarra - Return only WorkoutExercise that matters (with same id of the workout)
	for i1 := 0; i1 < len(workout.Exercises); i1++ {
		var exerciseSliceWithSameIdOfActualWorkout models.WorkoutExercise

		if workout.Exercises[i1].WorkoutExercises[0].WorkoutID == workout.Id {
			exerciseSliceWithSameIdOfActualWorkout = workout.Exercises[i1].WorkoutExercises[0]

			//Removing from slice the item we just added to another slice (the last one)
			// workout.Exercises[i1].WorkoutExercises = workout.Exercises[i1].WorkoutExercises[:len(workout.Exercises[i1].WorkoutExercises)-1]

		} else {
			fmt.Println("Something abnormal here. The first item of the workout.Exercises[i1].WorkoutExercises array shouldn't have a workoutID different from the actual workout.Id. Expected " + strconv.Itoa(workout.Id) + ", but got " + strconv.Itoa(workout.Exercises[i1].WorkoutExercises[0].WorkoutID))
			// Maybe this exercise is repeated in this week, and there is another one coming in the next days.
			// If its only one, the following gambiarra will solve.
			// But it should be refactored to be dynamic, or the query should Preload only exercises with WorkoutId Equal or Less than today
			if workout.Exercises[i1].WorkoutExercises[0].WorkoutID > workout.Id {
				//Remove workout.Exercises[i1].WorkoutExercises[0].WorkoutID from slice
				_, workout.Exercises[i1].WorkoutExercises = workout.Exercises[i1].WorkoutExercises[0], workout.Exercises[i1].WorkoutExercises[1:]

				if workout.Exercises[i1].WorkoutExercises[0].WorkoutID == workout.Id {
					fmt.Println("Mitigated successfully by popping the first item from slice, which has WorkoutID greaten that current WorkoutID")
					exerciseSliceWithSameIdOfActualWorkout = workout.Exercises[i1].WorkoutExercises[0]
				}
			}
		}

		var workoutExerciseSliceWithoutTheRecordsWithLoadValueEqualsToNull []models.WorkoutExercise
		secondItemIndex := 1
		for i2 := secondItemIndex; i2 < len(workout.Exercises[i1].WorkoutExercises); i2++ {

			if workout.Exercises[i1].WorkoutExercises[i2].WorkoutID == workout.Id {
				fmt.Println("Something abnormal here. The rest of the items of the workout.Exercises[i1].WorkoutExercises array shouldn't have a workoutID equals the actual workout.Id. Expected " + strconv.Itoa(workout.Id) + ", but got " + strconv.Itoa(workout.Exercises[i1].WorkoutExercises[i2].WorkoutID))
			} else {
				if workout.Exercises[i1].WorkoutExercises[i2].LoadValue == nil {
					fmt.Println("The history record with ExerciseID " + strconv.Itoa(workout.Exercises[i1].WorkoutExercises[i2].ExerciseID) + ", with WorkoutID " + strconv.Itoa(workout.Exercises[i1].WorkoutExercises[i2].WorkoutID) + ", with a null LoadValue will be omitted from response")
				} else {
					workoutExerciseSliceWithoutTheRecordsWithLoadValueEqualsToNull = append(workoutExerciseSliceWithoutTheRecordsWithLoadValueEqualsToNull, workout.Exercises[i1].WorkoutExercises[i2])
				}
				//It's history of exercises (exercise log)

				// break
				//TODO Limit the history to ten, but all then should be non null. This query will be good on the database, but can make the code complex because the actual exercise of this workout should be null... just the other ones that doesn't
			}
		}
		workout.Exercises[i1].WorkoutExercises = append([]models.WorkoutExercise{exerciseSliceWithSameIdOfActualWorkout}, workoutExerciseSliceWithoutTheRecordsWithLoadValueEqualsToNull...)
		//This "..." is from here https://stackoverflow.com/questions/16248241/concatenate-two-slices-in-go
	}

	//TODO: fix gambiarra - sort
	sort.Slice(workout.Exercises, func(i, j int) bool {
		return workout.Exercises[i].WorkoutExercises[0].Sort < workout.Exercises[j].WorkoutExercises[0].Sort
	})

	if workout.Id == 0 {
		// c.Writer.WriteHeader(http.StatusNoContent) //TODO: Studying, may be useful in queries, but I may be using query wrong, as pointed in a TODO comment above
		c.Writer.WriteHeader(http.StatusNotFound)
	} else {
		//TODO: use a real cache implementation (https://github.com/eko/gocache)
		dateWorkoutCache[dateWorkout] = workout

		//IndentedJSON is meant to be used on development only, because pretty print uses more CPU
		c.JSON(http.StatusOK, workout)
	}
}

func UpdateWorkoutExercise(c *gin.Context) {
	//TODO: create just once
	DB := db.Init()
	h := handlers.New(DB)

	var workoutExercise models.WorkoutExercise
	idWorkout, _ := strconv.Atoi(c.Params.ByName("idWorkout"))
	idExercise, _ := strconv.Atoi(c.Params.ByName("idExercise"))

	workoutExerciseFromDatabaseGet := h.GetWorkoutExercise(idWorkout, idExercise)

	if workoutExerciseFromDatabaseGet.WorkoutID == 0 || workoutExerciseFromDatabaseGet.ExerciseID == 0 {
		c.AbortWithStatus(404)
		return

	}

	// c.BindJSON(&workoutExercise)
	c.ShouldBindJSON(&workoutExercise)

	// workoutExerciseFromDatabaseUpdate := h.UpdateWorkoutExercise(workoutExercise)

	// c.JSON(http.StatusOK, workoutExerciseFromDatabaseUpdate)
	c.JSON(http.StatusOK, workoutExerciseFromDatabaseGet)
}
