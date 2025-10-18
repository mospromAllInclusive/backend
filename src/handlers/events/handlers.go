package events

import (
	"backend/src/handlers"
	"backend/src/modules/web_sockets"
	"backend/src/services"
)

func NewHandlers(
	usersService services.IUsersService,
	tablesHub *web_sockets.Hub,
) []handlers.IHandler {
	return []handlers.IHandler{
		newCellBusyHandler(tablesHub, usersService),
		newCellFreeHandler(tablesHub),
	}
}
