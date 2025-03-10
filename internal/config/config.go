package config

import (
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

	storeMasterPath := "storeMaster.csv"

	return &Config{
		Port:            port,
		StoreMasterPath: storeMasterPath,
		MaxWorkers:      10, 
	}
}
