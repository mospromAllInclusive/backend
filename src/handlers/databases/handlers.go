package databases

import (
	"backend/src/handlers"
	"backend/src/services"
)

func NewHandlers(
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
	usersService services.IUsersService,
) []handlers.IHandler {
	return []handlers.IHandler{
		newCreateDatabaseHandler(databasesService),
		newListDatabasesHandler(databasesService, tablesService),
		newGetDatabaseTablesHandler(tablesService, databasesService),
		newUsersHandler(databasesService),
		newSetRoleHandler(databasesService, usersService),
		newDeleteUserHandler(databasesService),
		newRoleHandler(databasesService),
	}
}
