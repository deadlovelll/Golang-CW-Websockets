package database

import (
	"log"

	"github.com/joho/godotenv"
)

// DatabaseConfig holds the database configuration values needed to establish a connection.
// It includes fields for database type, user credentials, database name, host, password, and SSL mode.
type DatabaseConfig struct {
	Type     string // Type of the database (e.g., "postgres")
	User     string // Database username
	Name     string // Name of the database
	Host     string // Host address of the database
	Password string // Password for the database user
	SslMode  string // SSL mode for the connection (e.g., "disable")
}

// LoadEnv loads environment variables from a .env file.
// It uses the godotenv package to read the file and logs a fatal error if the file cannot be loaded.
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
