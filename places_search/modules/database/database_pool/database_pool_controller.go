package database

import (
	"fmt"
	"places_search/modules/database/database"
	"sync"
)

// DatabasePoolController manages a singleton instance of the database connection.
// It provides methods to retrieve the database instance, initialize the connection,
// and gracefully shut down the connection pool.
type DatabasePoolController struct {
	// db holds the singleton instance of the Database.
	db *database.Database
	// mu ensures concurrent-safe access to the db instance.
	mu sync.Mutex
}

// GetDB returns the current database instance. If no instance exists, it creates one.
// This ensures that only one instance of the database connection is used throughout the application.
func (dpc *DatabasePoolController) GetDB() *database.Database {
	dpc.mu.Lock()
	defer dpc.mu.Unlock()

	if dpc.db == nil {
		dpc.db = database.GetDatabaseInstance()
	}
	return dpc.db
}

// Startup initializes the application by establishing the database connection.
// It retrieves the singleton instance using GetDB and then calls Connect to open the connection.
func (dpc *DatabasePoolController) Startup() {
	fmt.Println("Starting App...")
	db := dpc.GetDB()
	db.Connect() // Establish the database connection.
}

// Shutdown gracefully closes the database connection.
// It calls CloseAll to close all active connections, and recovers from any panic during shutdown.
func (dpc *DatabasePoolController) Shutdown() {
	fmt.Println("Shutting down app...")

	// Recover from any panic during shutdown to prevent application crash.
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Failed to close the database pool: %v. Application will stop. Traceback:\n%v\n", r, r)
		}
	}()

	if dpc.db != nil {
		dpc.db.CloseAll() // CloseAll is expected to close all connections.
		fmt.Println("Shutdown completed successfully.")
	}
}
