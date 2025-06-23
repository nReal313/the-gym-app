package handlers

import (
	"encoding/json"
	"net/http"
	"the-gym-app/internal/models"
	"the-gym-app/internal/services"
)

type WorkoutHandler struct {
	dbService *services.DatabaseService
}

func NewWorkoutHandler(dbService *services.DatabaseService) *WorkoutHandler {
	return &WorkoutHandler{dbService: dbService}
}

func (h *WorkoutHandler) LogWorkout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var workout models.WorkoutLog
	if err := json.NewDecoder(r.Body).Decode(&workout); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO : save the workout to a db
	if err := h.dbService.SaveWorkout(&workout); err != nil {
		http.Error(w, "Failed to save workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Workout logged successfully",
	})

}

func (h *WorkoutHandler) GetAllWorkouts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	logs, err := h.dbService.GetWorkouts()
	if err != nil {
		http.Error(w, "Unable to fetch workouts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
