package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Workout struct {
	Exercise string `json:"exercise"`
	Reps     int    `json:"reps"`
	ID       int    `json:"id"`
}

var workouts []Workout
var nextid = 0

func workoutsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(workouts)
	} else if r.Method == "POST" {
		var newWorkout Workout
		if err := json.NewDecoder(r.Body).Decode(&newWorkout); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		nextid = nextid + 1
		newWorkout.ID = nextid
		workouts = append(workouts, newWorkout)
		fmt.Fprintf(w, "New workout added")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func main() {
	http.HandleFunc("/workouts", workoutsHandler)

	fmt.Println("Server started on :8080")

	err := http.ListenAndServe(":8080", nil)
	fmt.Println("ListenAndServe error:", err)
}
