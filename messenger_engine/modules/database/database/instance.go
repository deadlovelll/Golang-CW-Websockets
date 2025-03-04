package database

import (
	"os"
	"sync"
)

var (
	// dbInstance holds the singleton instance of Database.
	dbInstance *Database
	// dbInstanceOnce ensures the singleton is initialized only once.
	dbInstanceOnce sync.Once
)

// GetDatabaseInstance returns a singleton instance of Database.
//
// It performs the following steps:
// 1. Loads the environment variables by calling LoadEnv().
// 2. Creates a new Database configuration using environment variables:
//    - DATABASE_TYPE
//    - DATABASE_USER
//    - DATABASE_PASSWORD
//    - DATABASE_NAME
//    - DATABASE_HOST
//    The SSL mode is hardcoded to "disable" (adjust if necessary).
// 3. Connects to the database by invoking the exported Connect() method on the Database instance.
// 4. Returns the singleton Database instance.
//
// Subsequent calls to this function will return the same instance.
func GetDatabaseInstance() *Database {
	dbInstanceOnce.Do(func() {
		LoadEnv()
		dbInstance = &Database{
			config: &DatabaseConfig{
				Type:     os.Getenv("DATABASE_TYPE"),
				User:     os.Getenv("DATABASE_USER"),
				Password: os.Getenv("DATABASE_PASSWORD"),
				Name:     os.Getenv("DATABASE_NAME"),
				Host:     os.Getenv("DATABASE_HOST"),
				SslMode:  "disable", // Adjust as necessary or load from an environment variable.
			},
		}
		dbInstance.Connect()
	})
	return dbInstance
}
