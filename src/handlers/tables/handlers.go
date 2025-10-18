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
	changelogService services.IChangelogService,
	tablesHub *web_sockets.Hub,
	usersHub *web_sockets.Hub,
) []handlers.IHandler {
	return []handlers.IHandler{
		newCreateTableHandler(usersHub, tablesService, databasesService),
		newImportTableHandler(usersHub, tablesService, databasesService, fileService),
		newAddColumnHandler(tablesHub, tablesService, databasesService, changelogService),
		newEditColumnHandler(tablesHub, tablesService, databasesService, changelogService),
		newDeleteColumnHandler(tablesHub, tablesService, databasesService, changelogService),
		newRestoreColumnHandler(tablesHub, tablesService, databasesService),
		newAddRowHandler(tablesHub, tablesService, databasesService, changelogService),
		newDeleteRowHandler(tablesHub, tablesService, databasesService, changelogService),
		newMoveRowHandler(tablesHub, tablesService, databasesService),
		newReadTableHandler(tablesService, databasesService),
		newExportTableHandler(tablesService, databasesService),
		newRestoreRowHandler(tablesHub, tablesService, databasesService),
		newSetCellValueHandler(tablesHub, tablesService, databasesService),
		newInfoHandler(tablesService, databasesService),
		newDeleteTableHandler(tablesHub, usersHub, tablesService, databasesService),
		newRestoreTableHandler(usersHub, tablesService, databasesService),
	}
}
