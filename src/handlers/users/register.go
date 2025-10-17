package users

import (
	"backend/src/handlers"
	"backend/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type registerHandler struct {
	authService services.IAuthService
	userService services.IUsersService
}

func newRegisterHandler(
	authService services.IAuthService,
	userService services.IUsersService,
) handlers.IHandler {
	return &registerHandler{
		authService: authService,
		userService: userService,
	}
}

func (h *registerHandler) Handle(c *gin.Context) {
	req := registerRequestDto{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	user, err := h.userService.AddUser(c, req.ToUser())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.authService.Login(c, req.Email, req.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newLoginResponse(token, user))
}

func (h *registerHandler) Path() string {
	return "/users/register"
}

func (h *registerHandler) Method() string {
	return http.MethodPost
}

func (h *registerHandler) AuthRequired() bool {
	return false
}
