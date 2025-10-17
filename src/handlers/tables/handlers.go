package tables

import (
	"backend/src/handlers"
	"backend/src/services"
)

func NewHandlers(
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
) []handlers.IHandler {
	return []handlers.IHandler{
		newCreateTableHandler(tablesService, databasesService),
		newAddColumnHandler(tablesService, databasesService),
		newAddRowHandler(tablesService, databasesService),
		newDeleteRowHandler(tablesService, databasesService),
		newMoveRowHandler(tablesService, databasesService),
		newReadTableHandler(tablesService, databasesService),
		newRestoreRowHandler(tablesService, databasesService),
		newSetCellValueHandler(tablesService, databasesService),
		newInfoHandler(tablesService, databasesService),
	}
}
