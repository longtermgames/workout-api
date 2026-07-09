package db

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

type Workout struct {
	Exercise  string `json:"exercise"`
	Reps      int    `json:"reps"`
	ID        int    `json:"id"`
	CreatedAt string `json:"created_at"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetWorkouts(userID int) ([]Workout, error) {
	rows, err := DB.Query("SELECT id, exercise, reps, created_at FROM workouts WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workouts := []Workout{}
	for rows.Next() {
		var w Workout
		if err := rows.Scan(&w.ID, &w.Exercise, &w.Reps, &w.CreatedAt); err != nil {
			return nil, err
		}
		workouts = append(workouts, w)
	}
	return workouts, nil
}

func CreateWorkout(exercise string, reps, userID int) error {
	_, err := DB.Exec(
		"INSERT INTO workouts (exercise, reps, user_id) VALUES ($1, $2, $3)",
		exercise,
		reps,
		userID,
	)

	return err
}

func GetWorkout(id, userID int) (Workout, error) {
	row := DB.QueryRow(
		"SELECT id, exercise, reps, created_at FROM workouts WHERE id = $1 and user_id = $2",
		id,
		userID,
	)
	var workout Workout
	err := row.Scan(
		&workout.ID,
		&workout.Exercise,
		&workout.Reps,
		&workout.CreatedAt,
	)
	return workout, err
}

func UpdateWorkout(id, userID int, exercise string, reps int) (int64, error) {
	result, err := DB.Exec(
		"UPDATE workouts SET exercise = $1, reps = $2 WHERE id = $3 and user_id = $4",
		exercise,
		reps,
		id,
		userID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
} // возвращает rowsAffected

func DeleteWorkout(id, userID int) (int64, error) {
	result, err := DB.Exec("DELETE FROM workouts where id = $1 and user_id = $2",
		id,
		userID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func CreateUser(username, hashedPassword string) error {
	_, err := DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, hashedPassword)
	return err
}

func GetUserByUsername(username string) (User, error) {
	var u User
	err := DB.QueryRow("SELECT id, password FROM users WHERE username = $1", username).Scan(&u.ID, &u.Password)
	return u, err
}

func Connect() error {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "user=naji dbname=workout_db sslmode=disable"
	}
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	return DB.Ping()
}
