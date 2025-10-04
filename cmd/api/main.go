package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/errol-vas/shiftplanner/internal/config"
	"github.com/errol-vas/shiftplanner/internal/handlers"
	"github.com/errol-vas/shiftplanner/internal/logger"
)

func main() {

	// Load Config
	cfg := config.Load()

	// Initialise Logger
	logger.Init()

	logger.Info(fmt.Sprintf("Starting application in %s mode", cfg.Env))

	// Set server to healthy
	atomic.StoreInt32(&handlers.Health, 1)

	// API Routes
	http.HandleFunc("/api/health", handlers.HealthCheck)
	http.HandleFunc("/api/version", handlers.Version(&cfg))

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start the Server
	logger.Info(fmt.Sprintf("Starting server on port %s", cfg.Port))
	log.Fatal(server.ListenAndServeTLS("server.crt", "server.key"))
}
