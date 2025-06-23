package main

import (
	"fmt"
	"log"
	"net/http"
	"the-gym-app/internal/handlers"
	"the-gym-app/internal/services"
)

func main() {
	//Initialise database
	dbService, err := services.NewDatabaseService()
	if err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}

	//Initialise handlers
	workoutHandler := handlers.NewWorkoutHandler(dbService)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to The Gym App")
	})

	http.HandleFunc("/api/workouts", workoutHandler.LogWorkout)

	http.HandleFunc("/api/workouts/findAll", workoutHandler.GetAllWorkouts)

	fmt.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
