package basecontroller

import (
	"places_search/modules/database/database"
)

// BaseController serves as a base structure for other controllers that require database access.
//
// Fields:
//   - Database: A pointer to the Database instance that provides access to database operations.
type BaseController struct {
	Database *database.Database
}
