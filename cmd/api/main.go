package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/errol-vas/shiftplanner/internal/handlers"
)

func main() {

	// Set server to healthy
	atomic.StoreInt32(&handlers.Health, 1)

	// API Routes
	http.HandleFunc("/api/health", handlers.HealthCheck)

	server := &http.Server{
		Addr:         ":5000",
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start the Server
	log.Println("Starting server on : 5000")
	log.Fatal(server.ListenAndServe())
}
