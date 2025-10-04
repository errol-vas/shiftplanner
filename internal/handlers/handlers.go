package handlers

import (
	"encoding/json"
	"net/http"
	"sync/atomic"

	"github.com/errol-vas/shiftplanner/internal/config"
	"github.com/errol-vas/shiftplanner/internal/constants"
)

var Health int32

type versionResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Env     string `json:"env"`
	Port    string `json:"port"`
}

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

func Version(cfg *config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := versionResponse{
			Name:    constants.AppName,
			Version: constants.AppVersion,
			Env:     cfg.Env,
			Port:    cfg.Port,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}
