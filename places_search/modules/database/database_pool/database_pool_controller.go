package database

import (
	"fmt"
	"places_search/modules/database/database"
	"sync"
)

// DatabasePoolController manages a singleton instance of the database connection.
type DatabasePoolController struct {
	db  *database.Database
	mu  sync.Mutex
}

// GetDB returns the current database instance, creating one if needed.
func (dpc *DatabasePoolController) GetDB() *database.Database {
	dpc.mu.Lock()
	defer dpc.mu.Unlock()

	if dpc.db == nil {
		dpc.db = database.GetDatabaseInstance()
	}
	return dpc.db
}

// Startup initializes the application by establishing the database connection.
func (dpc *DatabasePoolController) Startup() {
	fmt.Println("Starting App...")
	db := dpc.GetDB()
	db.Connect() // Establish the database connection.
}

// Shutdown gracefully closes the database connection.
func (dpc *DatabasePoolController) Shutdown() {
	fmt.Println("Shutting down app...")

	// Recover from any panic during shutdown.
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
