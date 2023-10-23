package main

import (
	"net/http"
	// "os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"vassopoli.com/gym-api/pkg/controllers"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:8082"},
		// AllowOrigins:     []string{os.Getenv("ALLOWED_ORIGIN")},
		AllowMethods:     []string{"PUT", "PATCH", "GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/workouts", controllers.GetWorkout)
	r.PUT("/workout-exercises/:idWorkout/:idExercise", controllers.UpdateWorkoutExercise)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	// r.Run(":80")
}
