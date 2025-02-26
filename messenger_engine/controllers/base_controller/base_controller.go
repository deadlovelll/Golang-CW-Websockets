package basecontroller

import (
	"messenger_engine/modules/database"
)

type BaseController struct {
	Database *database.Database
}
