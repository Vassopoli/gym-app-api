package models

import "time"

type Muscle struct {
	Id   int    `json:"id" gorm:"primaryKey;column:id"`
	Name string `json:"name" gorm:"column:name"`
}

func (Muscle) TableName() string {
	return "muscle"
}

type LoadType struct {
	Id   int    `json:"id" gorm:"primaryKey;column:id"`
	Name string `json:"name" gorm:"column:name"`
}

func (LoadType) TableName() string {
	return "load_type"
}

type WorkoutType struct {
	Id   int    `json:"id" gorm:"primaryKey;column:id"`
	Name string `json:"name" gorm:"column:name"`
}

func (WorkoutType) TableName() string {
	return "workout_type"
}

type AdvancedTechnique struct {
	Id   int    `json:"id" gorm:"primaryKey;column:id"`
	Name string `json:"name" gorm:"column:name"`
}

func (AdvancedTechnique) TableName() string {
	return "advanced_technique"
}

type Exercise struct {
	// gorm.Model //TODO: using this gorm model implies in creating deleted_at, and other stuff listed here https://gorm.io/docs/conventions.html
	Id               int               `json:"id" gorm:"primaryKey;column:id"`
	PrimaryName      string            `json:"primaryName" gorm:"column:primary_name"`
	IdYoutubeVideo   string            `json:"idYoutubeVideo" gorm:"column:id_youtube_video"`
	MuscleID         string            `json:"idMuscle" gorm:"column:id_muscle"`
	Muscle           Muscle            `json:"muscle"`
	WorkoutExercises []WorkoutExercise `json:"workoutExercises"`
	Workouts         []Workout         `json:"workouts" gorm:"many2many:workout_exercise;foreignKey:id;joinForeignKey:idExercise;References:id;joinReferences:idWorkout"`
}

func (Exercise) TableName() string {
	return "exercise"
}

type Workout struct {
	Id                   int         `json:"id" gorm:"primaryKey;column:id"`
	AerobicTimeInMinutes int         `json:"aerobicTimeInMinutes" gorm:"column:aerobic_time_in_minutes"`
	Date                 time.Time   `json:"date" gorm:"column:date"`
	Letter               string      `json:"letter" gorm:"column:letter"`
	WorkoutTypeID        int         `json:"idWorkoutType" gorm:"column:id_workout_type"`
	WorkoutType          WorkoutType `json:"workoutType"`
	Exercises            []Exercise  `json:"exercises" gorm:"many2many:workout_exercise;foreignKey:id;joinForeignKey:idWorkout;References:id;joinReferences:idExercise"`
}

func (Workout) TableName() string {
	return "workout"
}

type WorkoutExercise struct {
	WorkoutID           int               `json:"idWorkout" gorm:"primaryKey;column:id_workout"`
	ExerciseID          int               `json:"idExercise" gorm:"primaryKey;column:id_exercise"`
	Sort                int               `json:"sort" gorm:"column:sort"`
	LoadValue           *float64          `json:"loadValue" gorm:"column:load_value"`
	Executed            bool              `json:"executed" gorm:"column:executed"`
	Sets                int               `json:"sets" gorm:"column:sets"`
	Repetitions         int               `json:"repetitions" gorm:"column:repetitions"`
	RestSeconds         int               `json:"restSeconds" gorm:"column:rest_seconds"`
	LoadTypeID          int               `json:"idLoadType" gorm:"column:id_load_type"`
	LoadType            LoadType          `json:"loadType"`
	AdvancedTechniqueID *int              `json:"idAdvancedTechnique" gorm:"column:id_advanced_technique"`
	AdvancedTechnique   AdvancedTechnique `json:"advancedTechnique"`
}

func (WorkoutExercise) TableName() string {
	return "workout_exercise"
}
