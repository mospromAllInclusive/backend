package databases

import (
	"backend/src/domains/entities"
	"backend/src/handlers"
	"backend/src/handlers/common"
	"backend/src/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type getDatabaseTablesHandler struct {
	tablesService    services.ITablesService
	databasesService services.IDatabasesService
}

func newGetDatabaseTablesHandler(
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
) handlers.IHandler {
	return &getDatabaseTablesHandler{
		tablesService:    tablesService,
		databasesService: databasesService,
	}
}

func (h *getDatabaseTablesHandler) Handle(c *gin.Context) {
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

	tables, err := h.tablesService.ListByDatabaseID(c, dbIDInt)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, common.NewTablesResponse(tables))
}

func (h *getDatabaseTablesHandler) Path() string {
	return "/databases/:id/tables"
}

func (h *getDatabaseTablesHandler) Method() string {
	return http.MethodGet
}

func (h *getDatabaseTablesHandler) AuthRequired() bool {
	return true
}
