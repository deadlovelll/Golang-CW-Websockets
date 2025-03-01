package goenv

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from a .env file, if available.
// Logs a warning if the file is missing, but does not terminate execution.
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found. Falling back to system environment variables.")
	}
}

// GetEnv retrieves the value of an environment variable.
// If the variable is not set, it returns the provided fallback default value.
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
