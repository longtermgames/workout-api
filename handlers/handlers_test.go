package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Invalid JSON body on POST /workouts requires auth first, so without a token
// it returns 401 Unauthorized before parsing the body.
func TestWorkoutsHandler_NoToken(t *testing.T) {
	req := httptest.NewRequest("POST", "/workouts", bytes.NewBufferString("это не json"))
	w := httptest.NewRecorder()

	WorkoutsHandler(w, req)

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
