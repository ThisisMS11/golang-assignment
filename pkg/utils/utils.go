package utils

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

// GenerateRandomID generates a random ID string
func GenerateRandomID(length int) string {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		log.Printf("Error generating random ID: %v", err)
		return ""
	}
	return hex.EncodeToString(bytes)
}