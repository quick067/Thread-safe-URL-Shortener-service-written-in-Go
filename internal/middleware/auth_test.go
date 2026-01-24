package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


func TestAuthMiddleware_BadRequest(t *testing.T) {
	var secretKey string = "testKey"
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {w.WriteHeader(http.StatusOK)})

	req := httptest.NewRequest("GET", "/", nil)

	handleToTest := AuthMiddleware(nextHandler, secretKey)
	recorder := httptest.NewRecorder()
	handleToTest.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("Expected: %d, got: %d", http.StatusUnauthorized, recorder.Code)
	}
}

func TestAuthMiddleware_Success(t *testing.T) {
	var secretKey string = "testKey"

	claims := jwt.MapClaims{"user_id": 1.0, "exp": time.Now().Add(time.Hour * 2).Unix()}
	testToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := testToken.SignedString([]byte(secretKey))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("user_id")
		if userID == nil {
			t.Errorf("UserID not found in ctx")
		}
		w.WriteHeader(http.StatusOK)
	})
	handlerToTest := AuthMiddleware(nextHandler, secretKey)
	recorder := httptest.NewRecorder()

	handlerToTest.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK{
		t.Errorf("Expected: %d, got: %d", http.StatusOK, recorder.Code)
	}
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	var secretKey string = "testKey"


	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer ")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handlerToTest := AuthMiddleware(nextHandler, secretKey)
	recorder := httptest.NewRecorder()

	handlerToTest.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized{
		t.Errorf("Expected: %d, got: %d", http.StatusUnauthorized, recorder.Code)
	}
}

func TestAuthMiddleware_WrongSignature(t *testing.T) {
	var secretKey string = "testKey"

	claims := jwt.MapClaims{"user_id": 1.0, "exp": time.Now().Add(time.Hour * 2).Unix()}
	testToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := testToken.SignedString([]byte("hacketSign"))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handlerToTest := AuthMiddleware(nextHandler, secretKey)
	recorder := httptest.NewRecorder()

	handlerToTest.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized{
		t.Errorf("Expected: %d, got: %d", http.StatusUnauthorized, recorder.Code)
	}
}