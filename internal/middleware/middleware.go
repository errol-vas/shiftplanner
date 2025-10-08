package middleware

import (
	"net/http"
	"time"

	"github.com/errol-vas/shiftplanner/internal/logger"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		logger.Info(
			r.Method + " " + r.URL.Path +
				" - Status: " + http.StatusText(rw.statusCode) +
				" (" + duration.String() + ")",
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
