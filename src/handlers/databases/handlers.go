package databases

import (
	"backend/src/handlers"
	"backend/src/modules/web_sockets"
	"backend/src/services"
)

func NewHandlers(
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
	usersService services.IUsersService,
	usersHub *web_sockets.Hub,
) []handlers.IHandler {
	return []handlers.IHandler{
		newCreateDatabaseHandler(usersHub, databasesService),
		newListDatabasesHandler(databasesService, tablesService),
		newGetDatabaseTablesHandler(tablesService, databasesService),
		newUsersHandler(databasesService),
		newSetRoleHandler(usersHub, databasesService, usersService),
		newDeleteUserHandler(usersHub, databasesService, tablesService),
		newRoleHandler(databasesService),
	}
}
