package users

import (
	"backend/src/handlers"
	"backend/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type loginHandler struct {
	authService services.IAuthService
}

func newLoginHandler(
	authService services.IAuthService,
) handlers.IHandler {
	return &loginHandler{
		authService: authService,
	}
}

func (h *loginHandler) Handle(c *gin.Context) {
	req := loginRequestDto{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	user, token, err := h.authService.Login(c, req.Email, req.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newLoginResponse(token, user))
}

func (h *loginHandler) Path() string {
	return "/users/login"
}

func (h *loginHandler) Method() string {
	return http.MethodPost
}

func (h *loginHandler) AuthRequired() bool {
	return false
}
