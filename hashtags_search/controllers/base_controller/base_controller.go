package basecontroller

import (
	"hashtags_search/modules/database/database"
)

type BaseController struct {
	Database *database.Database
}
