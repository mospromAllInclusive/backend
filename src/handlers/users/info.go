package users

import (
	"backend/src/handlers"
	"backend/src/services"
	"backend/src/services/users"
	"net/http"

	"github.com/gin-gonic/gin"
)

type infoHandler struct {
	userService services.IUsersService
}

func newInfoHandler(
	userService services.IUsersService,
) handlers.IHandler {
	return &infoHandler{
		userService: userService,
	}
}

func (h *infoHandler) Handle(c *gin.Context) {
	user, err := h.userService.FindUserByID(c, c.MustGet("user_id").(int64))
	if err != nil {
		if users.IsErrUserNotFound(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newUserInfoResponse(user))
}

func (h *infoHandler) Path() string {
	return "/users/info"
}

func (h *infoHandler) Method() string {
	return http.MethodGet
}

func (h *infoHandler) AuthRequired() bool {
	return true
}
