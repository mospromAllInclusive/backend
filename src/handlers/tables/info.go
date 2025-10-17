package tables

import (
	"backend/src/domains/entities"
	"backend/src/handlers"
	"backend/src/handlers/common"
	"backend/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type infoHandler struct {
	tablesService    services.ITablesService
	databasesService services.IDatabasesService
}

func newInfoHandler(
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
) handlers.IHandler {
	return &infoHandler{
		tablesService:    tablesService,
		databasesService: databasesService,
	}
}

func (h *infoHandler) Handle(c *gin.Context) {
	tableID := c.Param("id")
	if tableID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid table id"})
		return
	}

	table, err := h.tablesService.GetTableByID(c, tableID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if table == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "table not found"})
		return
	}

	authorized, err := h.databasesService.CheckUserRole(c, c.MustGet("user_id").(int64), table.DatabaseID, entities.RoleReader)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !authorized {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user does not have writer role"})
		return
	}

	c.JSON(http.StatusOK, common.NewTableResponse(table))
}

func (h *infoHandler) Path() string {
	return "/tables/:id/info"
}

func (h *infoHandler) Method() string {
	return http.MethodGet
}

func (h *infoHandler) AuthRequired() bool {
	return true
}
