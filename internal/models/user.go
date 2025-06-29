package models

import "time"

type User struct {
	ID              int       `json:"id" db:"id"`
	Username        string    `json:"username" db:"username"`
	Email           string    `json:"email" db:"email"`
	PasswordHash    string    `json:"-" db:"password_hash"`
	Role            string    `json:"role"`
	FitnessGoal     string    `json:"fitness_goal,omitempty" db:"fitness_goal"`
	ExperienceLevel string    `json:"experience_level,omitempty" db:"experience_level"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`

	Workouts []Workout `json:"workouts,omitempty"`
}
