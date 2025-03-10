package api

import (
	"net/http"
	"time"
	"kiranaclub/internal/config"
)

// SetupRouter configures and returns the HTTP router
func SetupRouter() *http.Server {
	cfg := config.LoadConfig()
	
	// Register handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/api/submit", SubmitJobHandler)
	mux.HandleFunc("/api/status", GetJobStatusHandler)

	// Create server with timeouts
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return server
}