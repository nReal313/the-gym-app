package models

import "time"

type WorkoutLog struct {
	ID        int       `json:"id" db:"id"`
	Exercise  string    `json:"exercise"`
	Sets      []Set     `json:"sets"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Set struct {
	ID        int     `json:"id" db:"id"`
	WorkoutID int     `json:"workout_id" db:"workout_id"`
	Reps      int     `json:"reps"`
	Weight    float64 `json:"weight"`
	RPE       float64 `json:"rpe"`
	SetNumber int     `json:"set_number" db:"set_number"`
}
