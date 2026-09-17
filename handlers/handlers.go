package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"workout-api/auth"
	"workout-api/db"
)

const (
	minUsername = 3
	maxUsername = 50
	minPassword = 6
	maxPassword = 72 // bcrypt ignores bytes past 72
	maxExercise = 100
	maxReps     = 10000
)

// serverError logs the real cause and returns a generic message so internal
// details (SQL, driver errors) never leak to the client.
func serverError(w http.ResponseWriter, context string, err error) {
	log.Printf("%s: %v", context, err)
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

// authUserID authenticates the request; on failure it writes 401 and returns false.
func authUserID(w http.ResponseWriter, r *http.Request) (int, bool) {
	userID, err := auth.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return 0, false
	}
	return userID, true
}

// workoutID parses the ?id= query param; on failure it writes 400 and returns false.
func workoutID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid or missing id", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var user db.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	user.Username = strings.TrimSpace(user.Username)
	if len(user.Username) < minUsername || len(user.Username) > maxUsername {
		http.Error(w, "Username must be 3-50 characters", http.StatusBadRequest)
		return
	}
	if len(user.Password) < minPassword || len(user.Password) > maxPassword {
		http.Error(w, "Password must be 6-72 characters", http.StatusBadRequest)
		return
	}

	hashed, err := auth.HashPassword(user.Password)
	if err != nil {
		serverError(w, "hash password", err)
		return
	}
	if err := db.CreateUser(user.Username, hashed); err != nil {
		if errors.Is(err, db.ErrUserExists) {
			http.Error(w, "Username already taken", http.StatusConflict)
			return
		}
		serverError(w, "create user", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered"))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var user db.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	storedUser, err := db.GetUserByUsername(user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Same message as a wrong password so an attacker can't tell which
			// usernames exist (avoids user enumeration).
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		serverError(w, "lookup user", err)
		return
	}
	if err := auth.CheckPassword(storedUser.Password, user.Password); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	tokenString, err := auth.GenerateToken(storedUser.ID)
	if err != nil {
		serverError(w, "generate token", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func ListWorkoutsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := authUserID(w, r)
	if !ok {
		return
	}
	workouts, err := db.GetWorkouts(userID)
	if err != nil {
		serverError(w, "list workouts", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workouts)
}

func CreateWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := authUserID(w, r)
	if !ok {
		return
	}
	var newWorkout db.Workout
	if err := json.NewDecoder(r.Body).Decode(&newWorkout); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	if msg, ok := validateWorkout(newWorkout); !ok {
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	if err := db.CreateWorkout(strings.TrimSpace(newWorkout.Exercise), newWorkout.Reps, userID); err != nil {
		serverError(w, "create workout", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("New workout added"))
}

func GetWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := authUserID(w, r)
	if !ok {
		return
	}
	id, ok := workoutID(w, r)
	if !ok {
		return
	}
	workout, err := db.GetWorkout(id, userID)
	if err != nil {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)
}

func UpdateWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := authUserID(w, r)
	if !ok {
		return
	}
	id, ok := workoutID(w, r)
	if !ok {
		return
	}
	var updated db.Workout
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	if msg, ok := validateWorkout(updated); !ok {
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	rows, err := db.UpdateWorkout(id, userID, strings.TrimSpace(updated.Exercise), updated.Reps)
	if err != nil {
		serverError(w, "update workout", err)
		return
	}
	if rows == 0 {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}
	w.Write([]byte("Workout updated"))
}

func DeleteWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := authUserID(w, r)
	if !ok {
		return
	}
	id, ok := workoutID(w, r)
	if !ok {
		return
	}
	rows, err := db.DeleteWorkout(id, userID)
	if err != nil {
		serverError(w, "delete workout", err)
		return
	}
	if rows == 0 {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}
	w.Write([]byte("Workout deleted"))
}

// validateWorkout returns an error message and false if the workout is invalid.
func validateWorkout(wk db.Workout) (string, bool) {
	exercise := strings.TrimSpace(wk.Exercise)
	if exercise == "" || len(exercise) > maxExercise {
		return "Exercise must be 1-100 characters", false
	}
	if wk.Reps <= 0 || wk.Reps > maxReps {
		return "Reps must be between 1 and 10000", false
	}
	return "", true
}
