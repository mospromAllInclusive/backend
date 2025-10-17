package users

import (
	"backend/src/handlers"
	"backend/src/services"
)

func NewHandlers(
	authService services.IAuthService,
	userService services.IUsersService,
) []handlers.IHandler {
	return []handlers.IHandler{
		newLoginHandler(authService),
		newRegisterHandler(authService, userService),
		newInfoHandler(userService),
	}
}
