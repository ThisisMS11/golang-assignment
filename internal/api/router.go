package api

import (
	"net/http"
	"time"
	"kiranaclub/internal/config"
)

func SetupRouter() *http.Server {
	cfg := config.LoadConfig()
	
	/* Registering handlers */
	mux := http.NewServeMux()
	mux.HandleFunc("/api/submit", SubmitJobHandler)
	mux.HandleFunc("/api/status", GetJobStatusHandler)
	mux.HandleFunc("/api/jobs", GetJobsInformation)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return server
}