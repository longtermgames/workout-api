package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// POST /workouts requires auth first, so without a token it returns 401
// before the body is ever parsed.
func TestCreateWorkoutHandler_NoToken(t *testing.T) {
	req := httptest.NewRequest("POST", "/workouts", bytes.NewBufferString("это не json"))
	w := httptest.NewRecorder()

	CreateWorkoutHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("ожидал 401, получил %d", w.Code)
	}
}

// Register only accepts POST; a GET must be rejected with 405.
func TestRegisterHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest("GET", "/register", nil)
	w := httptest.NewRecorder()

	RegisterHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("ожидал 405, получил %d", w.Code)
	}
}

// Login only accepts POST; a GET must be rejected with 405.
func TestLoginHandler_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest("GET", "/login", nil)
	w := httptest.NewRecorder()

	LoginHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("ожидал 405, получил %d", w.Code)
	}
}

// A malformed JSON body is rejected with 400 before touching the database.
func TestRegisterHandler_InvalidBody(t *testing.T) {
	req := httptest.NewRequest("POST", "/register", bytes.NewBufferString("{broken"))
	w := httptest.NewRecorder()

	RegisterHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("ожидал 400, получил %d", w.Code)
	}
}

// Too-short credentials are rejected with 400 (validation runs before any DB call).
func TestRegisterHandler_ShortCredentials(t *testing.T) {
	req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(`{"username":"ab","password":"123"}`))
	w := httptest.NewRecorder()

	RegisterHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("ожидал 400, получил %d", w.Code)
	}
}
