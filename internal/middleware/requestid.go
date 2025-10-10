package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// contextKey is a private type to avoid context key collisions.
// We're not using plain strings to prevent interference from other packages.
type contextKey string

const RequestIDKey contextKey = "requestID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a new UUID
		id := uuid.New().String()

		// add the request id to the context so downstream handlers can access it
		ctx := context.WithValue(r.Context(), RequestIDKey, id)

		w.Header().Set("X-Request-ID", id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


func GetRequestID(ctx context.Context) string {
	if v := ctx.Value(RequestIDKey); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return "unknown"
}