package basecontroller

import (
	"hashtags_search/modules/database"
)

type BaseController struct {
	Database *database.Database
}
