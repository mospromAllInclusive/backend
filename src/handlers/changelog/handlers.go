package changelog

import (
	"backend/src/handlers"
	"backend/src/services"
)

func NewHandlers(
	changelogService services.IChangelogService,
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
) []handlers.IHandler {
	return []handlers.IHandler{
		newCellHandler(changelogService, tablesService, databasesService),
		newTableHandler(changelogService, tablesService, databasesService),
	}
}
