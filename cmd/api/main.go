package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/errol-vas/shiftplanner/internal/config"
	"github.com/errol-vas/shiftplanner/internal/handlers"
)

func main() {

	// Load Config
	cfg := config.Load()
	log.Printf("Starting application in %s mode on port %s", cfg.Env, cfg.Port)

	// Set server to healthy
	atomic.StoreInt32(&handlers.Health, 1)

	// API Routes
	http.HandleFunc("/api/health", handlers.HealthCheck)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start the Server
	log.Printf("Starting server on port %s", cfg.Port)
	log.Fatal(server.ListenAndServe())
}
