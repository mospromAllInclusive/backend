package tables

import (
	"backend/src/handlers"
	"backend/src/modules/web_sockets"
	"backend/src/services"
)

func NewHandlers(
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
	fileService services.IFileService,
	hub *web_sockets.Hub,
) []handlers.IHandler {
	return []handlers.IHandler{
		newCreateTableHandler(hub, tablesService, databasesService),
		newImportTableHandler(hub, tablesService, databasesService, fileService),
		newAddColumnHandler(hub, tablesService, databasesService),
		newEditColumnHandler(hub, tablesService, databasesService),
		newDeleteColumnHandler(hub, tablesService, databasesService),
		newRestoreColumnHandler(hub, tablesService, databasesService),
		newAddRowHandler(hub, tablesService, databasesService),
		newDeleteRowHandler(hub, tablesService, databasesService),
		newMoveRowHandler(hub, tablesService, databasesService),
		newReadTableHandler(tablesService, databasesService),
		newExportTableHandler(tablesService, databasesService),
		newRestoreRowHandler(hub, tablesService, databasesService),
		newSetCellValueHandler(hub, tablesService, databasesService),
		newInfoHandler(tablesService, databasesService),
		newDeleteTableHandler(hub, tablesService, databasesService),
		newRestoreTableHandler(hub, tablesService, databasesService),
	}
}
