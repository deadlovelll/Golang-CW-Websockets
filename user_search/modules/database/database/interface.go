package database

import (
	"database/sql"
)

// DatabaseInterface defines the behavior for a database.
type DatabaseInterface interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	GetConnection() *sql.DB
}