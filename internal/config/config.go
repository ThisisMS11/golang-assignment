package config

import (
	"log"
	"os"
)

// Config holds application configuration
type Config struct {
	Port            string
	StoreMasterPath string
	MaxWorkers      int
}

// LoadConfig loads configuration from environment variables or defaults
func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	storeMasterPath := os.Getenv("STORE_MASTER_PATH")
	log.Printf("storeMasterPath : %s", storeMasterPath)
	if storeMasterPath == "" {
		storeMasterPath = "/home/mohit/Desktop/Visual-Studio-Code/GoLang/KiranaClub-Assignment/storeMaster.csv"
	}

	return &Config{
		Port:            port,
		StoreMasterPath: storeMasterPath,
		MaxWorkers:      10, // Default number of concurrent workers
	}
}
