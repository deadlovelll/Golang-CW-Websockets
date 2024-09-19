package basecontroller

import (
	"user_search/modules/database"
)

type BaseController struct {
	Database *database.Database
}
