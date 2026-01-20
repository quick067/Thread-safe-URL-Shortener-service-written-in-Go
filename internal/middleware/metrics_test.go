package middleware 

import (
	"testing"
	"net/http"
	"net/http/httptest"
)

func TestMetricsMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {w.WriteHeader(http.StatusBadRequest)})

	handlerToTest := MetricsMiddleware(nextHandler)
	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	handlerToTest.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected: %d, got: %d", http.StatusBadRequest, recorder.Code)
	}
}