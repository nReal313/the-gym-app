package services

import (
	"database/sql"
	"the-gym-app/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type DatabaseService struct {
	db *sql.DB
}

func NewDatabaseService() (*DatabaseService, error) {
	db, err := sql.Open("sqlite3", "./gym_app.db")
	if err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}
	return &DatabaseService{db: db}, nil
}

func createTables(db *sql.DB) error {
	//create workout table
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS workouts (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		exercise TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS sets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		workout_id INTEGER NOT NULL,
		reps INTEGER NOT NULL,
		weight REAL NOT NULL,
		rpe REAL NOT NULL,
		set_number INTEGER NOT NULL,
		FOREIGN KEY (workout_id) REFERENCES workouts (id)
	)
	`)

	return err
}

func (s *DatabaseService) SaveWorkout(workout *models.WorkoutLog) error {
	//Start a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//Insert workout
	result, err := tx.Exec("INSERT INTO workouts (exercise) VALUES (?)", workout.Exercise)
	if err != nil {
		return err
	}

	workoutID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	//Insert sets
	for i, set := range workout.Sets {
		_, err := tx.Exec(
			"INSERT INTO sets (workout_id, reps, weight, rpe, set_number) VALUES (?, ?, ?, ?, ?)",
			workoutID, set.Reps, set.Weight, set.RPE, i+1,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *DatabaseService) GetWorkouts() ([]models.WorkoutLog, error) {
	var totalLogs []models.WorkoutLog
	rows, err := s.db.Query("SELECT id, exercise, created_at FROM workouts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var log models.WorkoutLog
		err := rows.Scan(&log.ID, &log.Exercise, &log.CreatedAt)
		if err != nil {
			return nil, err
		}
		setRows, err := s.db.Query("SELECT reps, weight, rpe, set_number FROM sets WHERE workout_id = ?", log.ID)
		if err != nil {
			return nil, err
		}
		defer setRows.Close()
		for setRows.Next() {
			var set models.Set
			err := setRows.Scan(&set.Reps, &set.Weight, &set.RPE, &set.SetNumber)
			if err != nil {
				return nil, err
			}
			log.Sets = append(log.Sets, set)
		}
		totalLogs = append(totalLogs, log)
	}
	return totalLogs, nil
}
