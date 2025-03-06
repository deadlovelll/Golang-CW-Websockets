package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Database wraps a sql.DB connection pool along with its configuration.
// It provides methods to establish the connection, retrieve the underlying
// connection object, and gracefully close the connection pool.
type Database struct {
	db     *sql.DB
	config *DatabaseConfig
}

// Connect initializes the database connection using the configuration provided
// in the Database instance. It constructs a connection string from the config,
// opens the connection, and verifies it with a ping. If any step fails, the
// function logs the error and terminates the application.
func (d *Database) Connect() {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s",
		d.config.Host, d.config.User, d.config.Password, d.config.Name, d.config.SslMode)

	var err error
	d.db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Verify the connection with a ping.
	if err := d.db.Ping(); err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	log.Println("Database connection established successfully")
}

// GetConnection returns the underlying *sql.DB connection pool.
// This connection can be used to perform database operations.
func (d *Database) GetConnection() *sql.DB {
	return d.db
}

// CloseAll gracefully closes the database connection pool.
// It logs an error and terminates the application if closing fails.
func (d *Database) CloseAll() {
	if d.db != nil {
		if err := d.db.Close(); err != nil {
			log.Fatalf("Error closing the database: %v", err)
		} else {
			log.Println("Database connection closed successfully")
		}
	}
}

// ReleaseConnection is a convenience method that wraps CloseAll.
// It is provided for semantic clarity when releasing the database connection.
func (d *Database) ReleaseConnection() {
	d.CloseAll()
}

func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.Query(query, args...) // Encapsulates direct access to *sql.DB
}