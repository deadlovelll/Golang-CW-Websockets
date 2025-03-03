package basecontroller

import (
	"places_search/modules/database"
)

type BaseController struct {
	Database *database.Database
}
