package databases

import (
	"backend/src/domains/entities"
	"backend/src/handlers"
	"backend/src/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type usersHandler struct {
	databasesService services.IDatabasesService
}

func newUsersHandler(
	databasesService services.IDatabasesService,
) handlers.IHandler {
	return &usersHandler{
		databasesService: databasesService,
	}
}

func (h *usersHandler) Handle(c *gin.Context) {
	dbID := c.Param("id")
	dbIDInt, err := strconv.ParseInt(dbID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID: " + err.Error()})
		return
	}

	authorized, err := h.databasesService.CheckUserRole(c, c.MustGet("user_id").(int64), dbIDInt, entities.RoleReader)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !authorized {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user does not have reader role"})
		return
	}

	databasesUsers, err := h.databasesService.GetDatabasesUsers(c, dbIDInt)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newDatabaseUsersListResponse(databasesUsers))
}

func (h *usersHandler) Path() string {
	return "/databases/:id/users"
}

func (h *usersHandler) Method() string {
	return http.MethodGet
}

func (h *usersHandler) AuthRequired() bool {
	return true
}
