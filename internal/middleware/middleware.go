package middleware

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/errol-vas/shiftplanner/internal/logger"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		reqID := GetRequestID(r.Context())

		// Client Metadata
		ip := extractIP(r)
		userAgent := r.Header.Get("User-Agent")
		referer := r.Header.Get("Referer")

		logLine := "[RequestID: " + reqID + "] " +
			r.Method + " " + r.URL.Path +
			" - Status: " + http.StatusText(rw.statusCode) +
			" (" + duration.String() + ")" +
			" - IP: " + ip +
			" - User-Agent: " + userAgent
		if referer != "" {
			logLine += " - Referer: " + referer
		}
		logger.Info(logLine)
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

func extractIP(r *http.Request) string {

	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		parts := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(parts[0])
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}
