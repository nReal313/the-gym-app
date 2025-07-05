package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"the-gym-app/internal/handlers"
	"the-gym-app/internal/middleware"
	"the-gym-app/internal/services"
)

func main() {
	//capturing flag for db cleanup
	var cleanup bool
	flag.BoolVar(&cleanup, "cleanup", false, "Delete db structure and data")

	flag.Parse()

	//Initialise database
	dbService, err := services.NewDatabaseService()
	if err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}

	//database cleanup
	if cleanup {
		log.Println("Cleaning up database...")
		if err := dbService.Cleanup(); err != nil {
			log.Fatal("Failed to cleanup database : ", err)
			return
		}
	}

	//Initialise handlers
	workoutHandler := handlers.NewWorkoutHandler(dbService)
	loginHandler := handlers.NewLoginHandler(dbService)

	http.HandleFunc("/signup", loginHandler.Signup)

	http.HandleFunc("/login", loginHandler.Login)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to The Gym App")
	})

	//endpoint to log workouts
	http.Handle("/api/workouts", middleware.MiddlewareHandler(http.HandlerFunc(workoutHandler.LogWorkout)))

	//endpoint to find a user's entire history of workouts
	http.Handle("/api/workouts/findAll", middleware.MiddlewareHandler(http.HandlerFunc(workoutHandler.GetAllWorkouts)))

	//endpoint to find user's set rep maxes
	http.Handle("/api/workouts/setRep", middleware.MiddlewareHandler(http.HandlerFunc(workoutHandler.GetSetRepMax)))

	fmt.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
