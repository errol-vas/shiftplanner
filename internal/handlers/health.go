package handlers

import (
	"net/http"
	"sync/atomic"
)

var Health int32

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if atomic.LoadInt32(&Health) == 1 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"available"}`))
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte(`{"status":"unavailable"}`))
}
