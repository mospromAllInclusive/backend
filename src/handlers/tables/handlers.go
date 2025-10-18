package tables

import (
	"backend/src/handlers"
	"backend/src/services"
)

func NewHandlers(
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
	fileService services.IFileService,
) []handlers.IHandler {
	return []handlers.IHandler{
		newCreateTableHandler(tablesService, databasesService),
		newImportTableHandler(tablesService, databasesService, fileService),
		newAddColumnHandler(tablesService, databasesService),
		newEditColumnHandler(tablesService, databasesService),
		newDeleteColumnHandler(tablesService, databasesService),
		newRestoreColumnHandler(tablesService, databasesService),
		newAddRowHandler(tablesService, databasesService),
		newDeleteRowHandler(tablesService, databasesService),
		newMoveRowHandler(tablesService, databasesService),
		newReadTableHandler(tablesService, databasesService),
		newExportTableHandler(tablesService, databasesService),
		newRestoreRowHandler(tablesService, databasesService),
		newSetCellValueHandler(tablesService, databasesService),
		newInfoHandler(tablesService, databasesService),
		newDeleteTableHandler(tablesService, databasesService),
		newRestoreTableHandler(tablesService, databasesService),
	}
}
