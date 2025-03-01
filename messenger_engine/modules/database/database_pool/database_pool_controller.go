package databasepool

import (
	"fmt"
	"sync"

	"messenger_engine/modules/database/database"
)

// DatabasePoolController manages the database pool.
type DatabasePoolController struct {
	Db *database.Database
	mu sync.Mutex
}

// GetDb returns the database instance, initializing it if necessary.
func (dp *DatabasePoolController) GetDb() *database.Database {
	dp.mu.Lock()
	defer dp.mu.Unlock()

	if dp.Db == nil {
		dp.Db = database.GetDatabaseInstance()  // Access the function directly
	}
	return dp.Db
}

// StartupEvent initializes the database connection during startup.
func (dp *DatabasePoolController) StartupEvent() {
	fmt.Println("Starting App...")

	db := dp.GetDb() // No need to call dp.Db.GetDb()
	db.Connect()     // Call the method directly
}

// ShutdownEvent handles the cleanup process during app shutdown.
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
