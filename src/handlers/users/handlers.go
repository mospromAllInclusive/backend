package users

import (
	"backend/src/handlers"
	"backend/src/services"
)

func NewHandlers(
	authService services.IAuthService,
) []handlers.IHandler {
	return []handlers.IHandler{
		newLoginHandler(authService),
	}
}
