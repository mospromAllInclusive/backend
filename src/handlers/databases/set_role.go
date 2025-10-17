package databases

import (
	"backend/src/domains/entities"
	"backend/src/handlers"
	"backend/src/services"
	"backend/src/services/users"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type setRoleHandler struct {
	databasesService services.IDatabasesService
	usersService     services.IUsersService
}

func newSetRoleHandler(
	databasesService services.IDatabasesService,
	usersService services.IUsersService,
) handlers.IHandler {
	return &setRoleHandler{
		databasesService: databasesService,
		usersService:     usersService,
	}
}

func (h *setRoleHandler) Handle(c *gin.Context) {
	req := setRoleRequestDto{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	dbID := c.Param("id")
	dbIDInt, err := strconv.ParseInt(dbID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID: " + err.Error()})
		return
	}

	authorized, err := h.databasesService.CheckUserRole(c, c.MustGet("user_id").(int64), dbIDInt, entities.RoleAdmin)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !authorized {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user does not have admin role"})
		return
	}

	_, err = h.usersService.FindUserByID(c, req.UserID)
	if err != nil {
		if users.IsErrUserNotFound(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = h.databasesService.UpsertUsersDatabase(c, &entities.UsersDatabase{
		DatabaseID: dbIDInt,
		UserID:     req.UserID,
		Role:       entities.Role(req.Role),
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *setRoleHandler) Path() string {
	return "/databases/:id/set-role"
}

func (h *setRoleHandler) Method() string {
	return http.MethodPost
}

func (h *setRoleHandler) AuthRequired() bool {
	return true
}
