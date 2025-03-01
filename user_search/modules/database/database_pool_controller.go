package database

import (
	"fmt"
	"sync"

	_ "github.com/lib/pq"
)

type DatabasePoolController struct {
	Db *Database
	mu sync.Mutex
}

func (dp *DatabasePoolController) GetDb() *Database {
	dp.mu.Lock()
	defer dp.mu.Unlock()

	if dp.Db == nil {
		dp.Db = GetDatabaseInstance()
	}
	return dp.Db
}

func (dp *DatabasePoolController) StartupEvent() {

	fmt.Println("Starting App...")

	db := dp.GetDb()

	db.connect()
}

func (dp *DatabasePoolController) ShutdownEvent() {

	fmt.Println("Shutdown app...")

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Failed to close the database pool: %v. Application would be stopped. Full traceback below.", r)
		}
	}()

	if dp.Db != nil {
		fmt.Println("Shutdown made successfully")
		dp.Db.CloseAll() // Предполагается, что CloseAll() закрывает все соединения
	}
}
