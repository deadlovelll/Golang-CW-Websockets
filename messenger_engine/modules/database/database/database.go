package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq" // PostgreSQL driver

	"messenger_engine/models/database_config"
	"messenger_engine/utls/env"
)

// Database represents a database connection with configuration.
type Database struct {
	db     *sql.DB
	config databaseconfig.DatabaseConfig
}

var (
	dbInstance     *Database
	dbInstanceOnce sync.Once
)

// GetDatabaseInstance returns a singleton instance of Database.
// It ensures that only one instance of the Database is created.
func GetDatabaseInstance() *Database {
	dbInstanceOnce.Do(func() {
		// Load environment variables
		goenv.LoadEnv()

		// Initialize DatabaseConfig using environment variables
		dbInstance = &Database{
			config: databaseconfig.DatabaseConfig{
				Type:     goenv.GetEnv("DATABASE_TYPE", "postgres"),
				User:     goenv.GetEnv("DATABASE_USER", "default_user"),
				Password: goenv.GetEnv("DATABASE_PASSWORD", ""),
				Name:     goenv.GetEnv("DATABASE_NAME", "default_db"),
				Host:     goenv.GetEnv("DATABASE_HOST", "localhost"),
				SslMode:  goenv.GetEnv("DATABASE_SSL_MODE", "disable"),
			},
		}

		// Connect to the database
		if err := dbInstance.Connect(); err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
	})
	return dbInstance
}

// Connect initializes the database connection.
// It constructs the connection string and establishes a connection to the database.
func (d *Database) Connect() error {
	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=%s",
		d.config.Host, d.config.User, d.config.Password, d.config.Name, d.config.SslMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error connecting to the database: %w", err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("database ping failed: %w", err)
	}

	d.db = db
	log.Println("Database connection established successfully.")
	return nil
}

// GetConnection returns the active database connection.
// If the connection is not established, it returns nil.
func (d *Database) GetConnection() *sql.DB {
	return d.db
}

// Close closes the database connection.
// It ensures the database connection is safely closed when no longer needed.
func (d *Database) Close() {
	if d.db != nil {
		if err := d.db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			log.Println("Database connection closed successfully.")
		}
	}
}