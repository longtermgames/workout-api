package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

type Workout struct {
	Exercise string `json:"exercise"`
	Reps     int    `json:"reps"`
	ID       int    `json:"id"`
}

var db *sql.DB

func workoutsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		workouts := []Workout{}
		rows, err := db.Query("SELECT id, exercise, reps FROM workouts")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var id int
			var exercise string
			var reps int
			if err := rows.Scan(&id, &exercise, &reps); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			workouts = append(workouts, Workout{ID: id, Exercise: exercise, Reps: reps})
			fmt.Printf("ID: %d, Exercise: %s, Reps: %d\n", id, exercise, reps)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(workouts)
	} else if r.Method == "POST" {
		var newWorkout Workout
		if err := json.NewDecoder(r.Body).Decode(&newWorkout); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err := db.Exec(
			"INSERT INTO workouts (exercise, reps) VALUES ($1, $2)",
			newWorkout.Exercise,
			newWorkout.Reps,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "New workout added")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func workoutHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if r.Method == "DELETE" {
		result, err := db.Exec("DELETE FROM workouts WHERE id = $1", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if rowsAffected == 0 {
			http.Error(w, "Workout not found", http.StatusNotFound)
		} else {
			fmt.Fprintf(w, "Workout with ID %d deleted", id)
		}
		return
	} else if r.Method == "PUT" {
		var updatedWorkout Workout
		if err := json.NewDecoder(r.Body).Decode(&updatedWorkout); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		result, err := db.Exec(
			"UPDATE workouts SET exercise = $1, reps = $2 WHERE id = $3",
			updatedWorkout.Exercise,
			updatedWorkout.Reps,
			id,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if rowsAffected == 0 {
			http.Error(w, "Workout not found", http.StatusNotFound)
		} else {
			fmt.Fprintf(w, "Workout updated")
		}
		return
	}

	row := db.QueryRow(
		"SELECT id, exercise, reps FROM workouts WHERE id = $1",
		id,
	)
	var workout Workout
	if err := row.Scan(
		&workout.ID,
		&workout.Exercise,
		&workout.Reps,
	); err != nil {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}

	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(workout)
	}
}

func main() {
	connStr := "user=naji dbname=workout_db sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}
	http.HandleFunc("/workouts", workoutsHandler)
	http.HandleFunc("/workout", workoutHandler)
	fmt.Println("Server started on :8080")

	err = http.ListenAndServe(":8080", nil)
	fmt.Println("ListenAndServe error:", err)
}
