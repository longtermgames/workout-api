package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// jwtSecret is set once at startup via SetSecret. It is intentionally not read
// from the environment at package-init time: that would make the package fail
// to load in tests and hide a missing secret behind an empty key.
var jwtSecret []byte

// SetSecret installs the HMAC secret used to sign and verify tokens.
func SetSecret(secret []byte) {
	jwtSecret = secret
}

// GetUserID validates the Bearer token from the request and returns its user id.
func GetUserID(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return 0, fmt.Errorf("missing token")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// Reject any token not signed with HMAC (guards against alg confusion,
		// e.g. a forged "none" or RS256 token).
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid claims")
	}
	rawID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid claims")
	}
	return int(rawID), nil
}

// GenerateToken issues a 24h HS256 token for the given user id.
func GenerateToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(jwtSecret)
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
