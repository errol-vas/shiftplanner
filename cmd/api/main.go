package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/errol-vas/shiftplanner/internal/config"
	"github.com/errol-vas/shiftplanner/internal/handlers"
	"github.com/errol-vas/shiftplanner/internal/logger"
	"github.com/errol-vas/shiftplanner/internal/middleware"
)

func main() {

	// Load Config
	cfg := config.Load()

	// Initialise Logger
	logger.Init()

	logger.Info(fmt.Sprintf("Starting application in %s mode", cfg.Env))

	// Set server to healthy
	atomic.StoreInt32(&handlers.Health, 1)

	mux := http.NewServeMux()

	// API Routes
	mux.HandleFunc("/api/health", handlers.HealthCheck)
	mux.HandleFunc("/api/version", handlers.Version(&cfg))

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      middleware.RequestID(middleware.Logging(mux)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful Shutdown
	go func() {
		logger.Info(fmt.Sprintf("Starting server on port %s", cfg.Port))
		err := server.ListenAndServeTLS("server.crt", "server.key")
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP Server error: %s", err)
		}
	}()

	// Listen for OS signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	// Graceful shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error(fmt.Sprintf("Could not gracefully shutdown the server: %v\n", err))
	} else {
		logger.Info("Server shutdown gracefully")
	}

}
