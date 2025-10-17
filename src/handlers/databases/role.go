package databases

import (
	"backend/src/handlers"
	"backend/src/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type roleHandler struct {
	databasesService services.IDatabasesService
}

func newRoleHandler(
	databasesService services.IDatabasesService,
) handlers.IHandler {
	return &roleHandler{
		databasesService: databasesService,
	}
}

func (h *roleHandler) Handle(c *gin.Context) {
	dbID := c.Param("id")
	dbIDInt, err := strconv.ParseInt(dbID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID: " + err.Error()})
		return
	}

	role, err := h.databasesService.GetUsersDatabaseRole(c, c.MustGet("user_id").(int64), dbIDInt)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if role == "" {
		c.Status(http.StatusForbidden)
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": role})
}

func (h *roleHandler) Path() string {
	return "/databases/:id/role"
}

func (h *roleHandler) Method() string {
	return http.MethodGet
}

func (h *roleHandler) AuthRequired() bool {
	return true
}
