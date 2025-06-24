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
	//create workoutlog table
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS workouts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			workout_name TEXT NOT NULL, 
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	//create  table
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS exercises (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		exercise TEXT NOT NULL,
		workout_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		FOREIGN KEY (workout_id) REFERENCES workouts (id)
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS sets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		exercise_id INTEGER NOT NULL,
		reps INTEGER NOT NULL,
		weight REAL NOT NULL,
		rpe REAL NOT NULL,
		set_number INTEGER NOT NULL,
		FOREIGN KEY (exercise_id) REFERENCES exercises (id)
	)
	`)

	return err
}

func (s *DatabaseService) SaveWorkout(workout *models.Workout) error {
	//Start a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec("INSERT INTO workouts (workout_name) VALUES (?)", workout.Name)
	if err != nil {
		return err
	}
	workoutID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	for _, exercise := range workout.Exercises {
		//Insert exercises
		result, err = tx.Exec("INSERT INTO exercises (exercise, workout_id) VALUES (?, ?)", exercise.Exercise, workoutID)
		if err != nil {
			return err
		}

		exerciseID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		//Insert sets
		for i, set := range exercise.Sets {
			_, err := tx.Exec(
				"INSERT INTO sets (exercise_id, reps, weight, rpe, set_number) VALUES (?, ?, ?, ?, ?)",
				exerciseID, set.Reps, set.Weight, set.RPE, i+1,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (s *DatabaseService) GetWorkouts() ([]models.Workout, error) {
	var totalLogs []models.Workout

	rows, err := s.db.Query("SELECT id, workout_name, created_at FROM workouts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var workout models.Workout
		err := rows.Scan(&workout.ID, &workout.Name, &workout.CreatedAt)
		if err != nil {
			return nil, err
		}
		exerciseRows, err := s.db.Query("SELECT id, exercise FROM exercises WHERE workout_id = ?", workout.ID)
		if err != nil {
			return nil, err
		}
		for exerciseRows.Next() {
			var exerciseLog models.ExerciseLog

			err := exerciseRows.Scan(&exerciseLog.ID, &exerciseLog.Exercise)
			if err != nil {
				return nil, err
			}
			setRows, err := s.db.Query("SELECT reps, weight, rpe, set_number FROM sets WHERE exercise_id = ?", exerciseLog.ID)
			if err != nil {
				return nil, err
			}
			for setRows.Next() {
				var set models.Set
				err := setRows.Scan(&set.Reps, &set.Weight, &set.RPE, &set.SetNumber)
				if err != nil {
					return nil, err
				}

				exerciseLog.Sets = append(exerciseLog.Sets, set)
			}
			setRows.Close()

			workout.Exercises = append(workout.Exercises, exerciseLog)
		}
		exerciseRows.Close()
		totalLogs = append(totalLogs, workout)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return totalLogs, nil
}

// func (s *DatabaseService) GetWorkout(date time.Time) (models.WorkoutLog, error) {
// 	rows, err := s.db.Query("SELECT id, exercise, created_at FROM ")
// }
