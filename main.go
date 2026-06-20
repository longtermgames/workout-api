package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var jwtSecret = []byte(getEnv("JWT_SECRET", "supersecret"))

type Workout struct {
	ID       int    `json:"id"`
	Exercise string `json:"exercise"`
	Reps     int    `json:"reps"`
	UserID   int    `json:"user_id"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// --- Auth ---

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", u.Username, string(hash))
	if err != nil {
		http.Error(w, "Username already taken", http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User registered")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	var stored User
	err := db.QueryRow("SELECT id, password FROM users WHERE username = $1", u.Username).Scan(&stored.ID, &stored.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(stored.Password), []byte(u.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": stored.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenStr})
}

// --- Middleware ---

func getUserID(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return 0, fmt.Errorf("missing token")
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := int(claims["user_id"].(float64))
	return userID, nil
}

// --- Workouts ---

func workoutsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method == "GET" {
		rows, err := db.Query("SELECT id, exercise, reps FROM workouts WHERE user_id = $1", userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		workouts := []Workout{}
		for rows.Next() {
			var wo Workout
			if err := rows.Scan(&wo.ID, &wo.Exercise, &wo.Reps); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			workouts = append(workouts, wo)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(workouts)

	} else if r.Method == "POST" {
		var wo Workout
		if err := json.NewDecoder(r.Body).Decode(&wo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err := db.Exec(
			"INSERT INTO workouts (exercise, reps, user_id) VALUES ($1, $2, $3)",
			wo.Exercise, wo.Reps, userID,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "New workout added")

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func workoutHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if r.Method == "DELETE" {
		result, err := db.Exec("DELETE FROM workouts WHERE id = $1 AND user_id = $2", id, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rows, _ := result.RowsAffected()
		if rows == 0 {
			http.Error(w, "Workout not found", http.StatusNotFound)
		} else {
			fmt.Fprintf(w, "Workout deleted")
		}

	} else if r.Method == "PUT" {
		var wo Workout
		if err := json.NewDecoder(r.Body).Decode(&wo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		result, err := db.Exec(
			"UPDATE workouts SET exercise = $1, reps = $2 WHERE id = $3 AND user_id = $4",
			wo.Exercise, wo.Reps, id, userID,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rows, _ := result.RowsAffected()
		if rows == 0 {
			http.Error(w, "Workout not found", http.StatusNotFound)
		} else {
			fmt.Fprintf(w, "Workout updated")
		}

	} else if r.Method == "GET" {
		var wo Workout
		err := db.QueryRow(
			"SELECT id, exercise, reps FROM workouts WHERE id = $1 AND user_id = $2",
			id, userID,
		).Scan(&wo.ID, &wo.Exercise, &wo.Reps)
		if err != nil {
			http.Error(w, "Workout not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(wo)

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "user=naji dbname=workout_db sslmode=disable"
	}
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = db.Ping(); err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/workouts", workoutsHandler)
	http.HandleFunc("/workout", workoutHandler)

	fmt.Println("Server started on :8080")
	if err = http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("ListenAndServe error:", err)
	}
}
