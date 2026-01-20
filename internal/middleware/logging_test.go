package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestLogging(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {w.WriteHeader(http.StatusOK)})

	req := httptest.NewRequest("GET", "/", nil)
	handlerToTest := LogMiddleware(nextHandler)

	recorder := httptest.NewRecorder()
	handlerToTest.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK{
		t.Errorf("Expected: %d, got: %d", http.StatusOK, recorder.Code)
	}
}