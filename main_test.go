package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostWorkout_InvalidBody(t *testing.T) {
	req := httptest.NewRequest("POST", "/workouts", bytes.NewBufferString("это не json"))
	w := httptest.NewRecorder()

	workoutsHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("ожидал 400, получил %d", w.Code)
	}
}

func TestWorkoutHandler_InvalidID(t *testing.T) {
	req := httptest.NewRequest("GET", "/workout?id=abc", nil)
	w := httptest.NewRecorder()

	workoutHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("ожидал 400, получил %d", w.Code)
	}
}

func TestWorkoutsHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest("DELETE", "/workouts", nil)
	w := httptest.NewRecorder()

	workoutsHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("ожидал 405, получил %d", w.Code)
	}
}
