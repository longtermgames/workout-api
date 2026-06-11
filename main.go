package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Workout struct {
	Exercise string `json:"exercise"`
	Reps     int    `json:"reps"`
	ID       int    `json:"id"`
}

var workouts = []Workout{}
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
func workoutHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	getworkout := Workout{ID: -1}
	for _, workout := range workouts {
		if workout.ID == id {
			getworkout = workout
			break
		}
	}
	if getworkout.ID == -1 {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}
	if r.Method == "DELETE" {
		for i, workout := range workouts {
			if workout.ID == id {
				workouts = append(workouts[:i], workouts[i+1:]...)
				fmt.Fprintf(w, "Workout deleted")
				return
			}
		}
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getworkout)
}

func main() {
	http.HandleFunc("/workouts", workoutsHandler)
	http.HandleFunc("/workout", workoutHandler)
	fmt.Println("Server started on :8080")

	err := http.ListenAndServe(":8080", nil)
	fmt.Println("ListenAndServe error:", err)
}
