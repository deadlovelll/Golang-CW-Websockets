package basecontroller

import (
	"log"

	"messenger_engine/modules/database"
)

// BaseController provides a common structure for controllers with database access.
type BaseController struct {
	db *database.Database
}

// NewBaseController initializes a new BaseController.
func NewBaseController(db *database.Database) *BaseController {
	if db == nil {
		log.Fatal("BaseController: database instance cannot be nil")
	}
	return &BaseController{db: db}
}

// GetDatabase returns the database instance.
func (bc *BaseController) GetDatabase() *database.Database {
	return bc.db
}
