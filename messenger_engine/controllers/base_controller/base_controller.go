package basecontroller

import (
	"log"

	"messenger_engine/modules/database"
)

// BaseController provides a common structure for controllers with database access.
type BaseController struct {
	Database *database.Database
}

// NewBaseController initializes a new BaseController.
func NewBaseController(Database *database.Database) *BaseController {
	if Database == nil {
		log.Fatal("BaseController: database instance cannot be nil")
	}
	return &BaseController{Database: Database}
}

// GetDatabase returns the database instance.
func (bc *BaseController) GetDatabase() *database.Database {
	return bc.Database
}
