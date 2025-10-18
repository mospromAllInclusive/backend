package tables

import (
	"backend/src/domains/entities"
	"backend/src/handlers"
	"backend/src/services"
	"backend/src/services/tables"
	"net/http"

	"github.com/gin-gonic/gin"
)

type readTableHandler struct {
	tablesService    services.ITablesService
	databasesService services.IDatabasesService
}

func newReadTableHandler(
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
) handlers.IHandler {
	return &readTableHandler{
		tablesService:    tablesService,
		databasesService: databasesService,
	}
}

func (h *readTableHandler) Handle(c *gin.Context) {
	var q entities.ReadTableParams
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tableID := c.Param("id")
	if tableID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid table id"})
		return
	}

	table, err := h.tablesService.GetTableByID(c, tableID, false)
	if err != nil {
		if tables.IsErrTableNotFound(err) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "table not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !table.ValidateParams(&q) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid filter"})
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

	rows, err := h.tablesService.ReadTable(c, table, q)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	total, err := h.tablesService.GetTotalRows(c, table, q)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newTableWithDataResponse(table, rows, total))
}

func (h *readTableHandler) Path() string {
	return "/tables/:id"
}

func (h *readTableHandler) Method() string {
	return http.MethodGet
}

func (h *readTableHandler) AuthRequired() bool {
	return true
}
