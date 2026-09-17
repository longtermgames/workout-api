package auth

import (
	"net/http"
	"testing"
)

func TestTokenRoundTrip(t *testing.T) {
	SetSecret([]byte("test-secret"))

	tok, err := GenerateToken(42)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	req := httptest_newRequestWithBearer(tok)
	id, err := GetUserID(req)
	if err != nil {
		t.Fatalf("get user id: %v", err)
	}
	if id != 42 {
		t.Errorf("ожидал 42, получил %d", id)
	}
}

func TestGetUserID_NoHeader(t *testing.T) {
	SetSecret([]byte("test-secret"))
	req, _ := http.NewRequest("GET", "/", nil)
	if _, err := GetUserID(req); err == nil {
		t.Error("ожидал ошибку при отсутствии заголовка Authorization")
	}
}

// A token signed with a different secret must be rejected (signature check works).
func TestGetUserID_WrongSecret(t *testing.T) {
	SetSecret([]byte("secret-A"))
	tok, err := GenerateToken(7)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	SetSecret([]byte("secret-B"))
	req := httptest_newRequestWithBearer(tok)
	if _, err := GetUserID(req); err == nil {
		t.Error("ожидал ошибку для токена, подписанного другим секретом")
	}
}

func httptest_newRequestWithBearer(token string) *http.Request {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}
