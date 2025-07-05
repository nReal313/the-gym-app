package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"the-gym-app/internal/middleware"
	"the-gym-app/internal/models"
	"the-gym-app/internal/services"

	"github.com/golang-jwt/jwt/v5"
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

	var workout models.Workout
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

func (h *WorkoutHandler) GetSetRepMax(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== GetSetRepMax handler called ===")
	log.Printf("URL: %s", r.URL.String())
	log.Printf("Method: %s", r.Method)

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	exercise := r.URL.Query().Get("exercise")
	if exercise == "" {
		http.Error(w, "exercise parameter is required", http.StatusBadRequest)
		return
	}

	repsStr := r.URL.Query().Get("reps")
	if repsStr == "" {
		http.Error(w, "reps parameter is required", http.StatusBadRequest)
		return
	}
	repsInt, err := strconv.Atoi(repsStr)
	if err != nil {
		http.Error(w, "Error during type conversion", http.StatusInternalServerError)
		return
	}

	if repsInt <= 0 {
		http.Error(w, "reps must be greater than 0", http.StatusBadRequest)
		return
	}

	//get the user from context
	log.Printf("Attempting to extract user from context...")
	userValue := r.Context().Value(middleware.ContextKey("user"))
	log.Printf("User value from context: %+v", userValue)

	userClaims, ok := userValue.(jwt.MapClaims)
	if !ok {
		log.Printf("Failed to convert user value to jwt.MapClaims")
		http.Error(w, "Invalid user context", http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully extracted user claims: %+v", userClaims)

	userName, ok := userClaims["username"].(string)
	if !ok {
		http.Error(w, "Invalid username : not a string!", http.StatusInternalServerError)
		return
	}

	userID, err := h.dbService.GetUserIdFromUsername(userName)
	if err != nil {
		log.Printf("Error fetching user ID: %v", err)
		if strings.Contains(err.Error(), "user not found") {
			http.Error(w, "User not found in database", http.StatusUnauthorized)
		} else {
			http.Error(w, "Error while fetching User ID", http.StatusInternalServerError)
		}
		return
	}
	setRep, err := h.dbService.GetSetRep(userID, exercise, repsInt)
	if err != nil {
		http.Error(w, "Error fetching set rep data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(setRep)
}
