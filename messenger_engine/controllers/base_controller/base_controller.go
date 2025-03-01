package basecontroller

import (
	"log"

	"messenger_engine/modules/database/database"
)

// BaseController provides a common structure for controllers that require database access.
type BaseController struct {
	Database *database.Database // Database instance used by the controller.
}

// NewBaseController initializes a new BaseController.
//
// This function ensures that a valid database instance is provided.
// If a nil database is passed, it logs a fatal error and stops execution.
//
// Parameters:
//   - Database: A pointer to a database instance.
//
// Returns:
//   - A pointer to a newly created BaseController instance.
func NewBaseController(Database *database.Database) *BaseController {
	if Database == nil {
		log.Fatal("BaseController: database instance cannot be nil")
	}
	return &BaseController{Database: Database}
}

// GetDatabase returns the associated database instance of the BaseController.
//
// Returns:
//   - A pointer to the Database instance.
func (bc *BaseController) GetDatabase() *database.Database {
	return bc.Database
}