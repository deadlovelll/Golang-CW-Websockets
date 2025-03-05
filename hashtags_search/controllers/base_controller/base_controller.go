package basecontroller

import (
	"hashtags_search/modules/database/database"
)

// BaseController provides a base structure for controllers that need database access.
type BaseController struct {
	// Database holds a reference to the database connection.
	Database database.DatabaseInterface
}