package users

import (
	"backend/src/handlers"
	"backend/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type listHandler struct {
	userService services.IUsersService
}

func newListHandler(
	userService services.IUsersService,
) handlers.IHandler {
	return &listHandler{
		userService: userService,
	}
}

func (h *listHandler) Handle(c *gin.Context) {
	usersList, err := h.userService.ListUsers(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newUsersListResponse(usersList))
}

func (h *listHandler) Path() string {
	return "/users/list"
}

func (h *listHandler) Method() string {
	return http.MethodGet
}

func (h *listHandler) AuthRequired() bool {
	return true
}
