package models

import "time"

type Workout struct {
	ID        int           `json:"id" db:"id"`
	Name      string        `json:"name" db:"workout_name"`
	Exercises []ExerciseLog `json:"exercises" db:"exercises"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
	UserID    int           `json:"user_id" db:"user_id"`
}

type ExerciseLog struct {
	ID        int       `json:"id" db:"id"`
	Exercise  string    `json:"exercise"`
	Sets      []Set     `json:"sets"`
	WorkoutID int       `json:"workout_id" db:"workout_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Set struct {
	ID         int     `json:"id" db:"id"`
	ExerciseID int     `json:"exercise_id" db:"exercise_id"`
	Reps       int     `json:"reps"`
	Weight     float64 `json:"weight"`
	RPE        float64 `json:"rpe"`
	SetNumber  int     `json:"set_number" db:"set_number"`
	// TODO : add remarks field here and make it optional
}

type SetRep struct {
	ExerciseName string  `json:"exercise_name"`
	Reps         int     `json:"reps"`
	Weight       float64 `json:"weight"`
}
