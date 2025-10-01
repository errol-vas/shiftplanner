package main

import (
	"log"
	"net/http"

	"github.com/errol-vas/shiftplanner/internal/handlers"
)




func main() {

	// API Routes
	http.HandleFunc("/api/health", handlers.HealthCheck)

	// Start the Server
	log.Println("Starting server on : 5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
