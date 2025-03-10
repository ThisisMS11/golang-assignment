package main

import (
	"log"
	"math/rand"
	"time"
	"kiranaclub/internal/api"
	"kiranaclub/internal/config"
	"kiranaclub/internal/store"
)

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Load configuration
	cfg := config.LoadConfig()

	// Load store master data
	if err := store.LoadStoreMaster(cfg.StoreMasterPath); err != nil {
		log.Fatalf("Failed to load store master: %v", err)
	}
	log.Printf("Loaded %d stores from store master", store.Count())

	// Set up and start HTTP server
	router := api.SetupRouter()
	port := cfg.Port

	log.Printf("Server starting on port %s...", port)
	if err := router.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}