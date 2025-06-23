package main

import (
	"fmt"
	"log"
	"net/http"
	"the-gym-app/internal/handlers"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to The Gym App")
	})

	http.HandleFunc("/api/workouts", handlers.LogWorkout)

	fmt.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
