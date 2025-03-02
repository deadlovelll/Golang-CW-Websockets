package basecontroller

import (
	"user_search/modules/database/database"
)

type BaseController struct {
	Database *database.Database
}
