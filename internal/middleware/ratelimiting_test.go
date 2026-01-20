package middleware

import(
	"testing"
	"net/http"
	"net/http/httptest"

	"golang.org/x/time/rate"
)

func TestRateLimit_Block(t *testing.T){
	rl := &RateLimiter{
		LimitMap: map[string]*rate.Limiter{},
	}
	nextHandler := http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {w.WriteHeader(http.StatusOK)})
	handlerToTest := rl.RateLimitMiddleware(nextHandler)
	for i := 0; i < 4; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "192.168.0.1:54321"
		recorder := httptest.NewRecorder()

		handlerToTest.ServeHTTP(recorder, req)

		if i > 2 && recorder.Code != http.StatusTooManyRequests{
			t.Errorf("Expected: %d, got: %d", http.StatusTooManyRequests, recorder.Code)
		}
		if i != 3 && recorder.Code != http.StatusOK{
			t.Errorf("Expected: %d, got: %d", http.StatusOK, recorder.Code)
		}
	}
}