package middleware

import (
	"net"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	LimitMap map[string]*rate.Limiter
	mtx      sync.Mutex
}

func (l *RateLimiter) getLimiter(ip string) *rate.Limiter {
	l.mtx.Lock()
	limiter, ok := l.LimitMap[ip]
	if !ok {
		limiter = rate.NewLimiter(1, 3)
		l.LimitMap[ip] = limiter
	}
	l.mtx.Unlock()
	return limiter
}

func (l *RateLimiter) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIp, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			userIp = r.RemoteAddr
		}
		userLimiter := l.getLimiter(userIp)
		if !userLimiter.Allow() {
			http.Error(w, "Rate limit exhausted", http.StatusTooManyRequests)
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
