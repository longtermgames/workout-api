package main

import (
	"fmt"
	"net/http"

	"workout-api/db"
	"workout-api/handlers"
)

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func main() {
	if err := db.Connect(); err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/register", corsMiddleware(handlers.RegisterHandler))
	http.HandleFunc("/login", corsMiddleware(handlers.LoginHandler))
	http.HandleFunc("/workouts", corsMiddleware(handlers.WorkoutsHandler))
	http.HandleFunc("/workout", corsMiddleware(handlers.WorkoutHandler))

	fmt.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("ListenAndServe error:", err)
	}
}
