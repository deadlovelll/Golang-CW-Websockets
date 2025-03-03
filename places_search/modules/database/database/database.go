package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Database wraps a sql.DB connection pool along with its configuration.
type Database struct {
	db     *sql.DB
	config *DatabaseConfig
}

// connect initializes the database connection using the provided configuration.
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

// GetConnection returns the underlying *sql.DB connection.
func (d *Database) GetConnection() *sql.DB {
	return d.db
}

// CloseAll gracefully closes the database connection pool.
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
func (d *Database) ReleaseConnection() {
	d.CloseAll()
}
