package databasepool

import (
	"fmt"
	"sync"

	"messenger_engine/modules/database/database"
)

// DatabasePoolController manages the database pool and ensures
// that only a single instance of the database connection is used.
type DatabasePoolController struct {
	Db *database.Database
	mu sync.Mutex
}

// GetDb returns the database instance, initializing it if necessary.
// It ensures thread safety by using a mutex lock.
func (dp *DatabasePoolController) GetDb() *database.Database {
	dp.mu.Lock()
	defer dp.mu.Unlock()

	if dp.Db == nil {
		dp.Db = database.GetDatabaseInstance() // Access the function directly
	}
	return dp.Db
}

// StartupEvent initializes the database connection during application startup.
// It ensures that the database instance is created and connected.
func (dp *DatabasePoolController) StartupEvent() {
	fmt.Println("Starting App...")

	db := dp.GetDb() // Retrieve the database instance
	db.Connect()     // Establish the database connection
}

// ShutdownEvent handles the cleanup process during application shutdown.
// It ensures that the database connection is properly closed to prevent resource leaks.
func (dp *DatabasePoolController) ShutdownEvent() {
	fmt.Println("Shutdown app...")

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Failed to close the database pool: %v. Application would be stopped. Full traceback below.\n", r)
		}
	}()

	if dp.Db != nil {
		fmt.Println("Shutdown made successfully")
		dp.Db.Close() // Close the database connection
	}
}