package middleware

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	operationsProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "total_processed_operations",
		Help: "The total number of processed events",
	}, []string{"method", "path", "status"})
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int 
}

func (sr *statusRecorder) WriteHeader (code int) {
	sr.statusCode = code
	sr.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware (next http.Handler) http.Handler {
	return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{
			ResponseWriter: w,
			statusCode: 200,
		}

		next.ServeHTTP(rec, r)

		operationsProcessed.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(rec.statusCode)).Inc()
	})
}