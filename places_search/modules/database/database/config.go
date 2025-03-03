package database

import (
	"log"

	"github.com/joho/godotenv"
)

// DatabaseConfig holds the database configuration values.
type DatabaseConfig struct {
	Type     string
	User     string
	Name     string
	Host     string
	Password string
	SslMode  string
}

// LoadEnv loads environment variables from a .env file.
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
