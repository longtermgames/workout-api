package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"workout-api/auth"
	"workout-api/db"
	"workout-api/handlers"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is required")
	}
	auth.SetSecret([]byte(secret))

	if err := db.Connect(); err != nil {
		log.Fatal("db connect: ", err)
	}

	// Go 1.22+ method-based routing: the pattern encodes the HTTP method, so the
	// mux returns 405 automatically and each handler does one thing.
	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", handlers.RegisterHandler)
	mux.HandleFunc("POST /login", handlers.LoginHandler)
	mux.HandleFunc("GET /workouts", handlers.ListWorkoutsHandler)
	mux.HandleFunc("POST /workouts", handlers.CreateWorkoutHandler)
	mux.HandleFunc("GET /workout", handlers.GetWorkoutHandler)
	mux.HandleFunc("PUT /workout", handlers.UpdateWorkoutHandler)
	mux.HandleFunc("DELETE /workout", handlers.DeleteWorkoutHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Timeouts protect against slow/idle clients (e.g. slowloris) that would
	// otherwise tie up connections forever.
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           corsMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("Server started on :" + port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
