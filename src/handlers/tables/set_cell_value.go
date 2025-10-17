package tables

import (
	"backend/src/domains/entities"
	"backend/src/handlers"
	"backend/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type setCellValueHandler struct {
	tablesService    services.ITablesService
	databasesService services.IDatabasesService
}

func newSetCellValueHandler(
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
) handlers.IHandler {
	return &setCellValueHandler{
		tablesService:    tablesService,
		databasesService: databasesService,
	}
}

func (h *setCellValueHandler) Handle(c *gin.Context) {
	req := setCellValueRequestDto{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

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

	authorized, err := h.databasesService.CheckUserRole(c, c.MustGet("user_id").(int64), table.DatabaseID, entities.RoleWriter)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !authorized {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user does not have writer role"})
		return
	}

	err = h.tablesService.SetCellValue(c, table.ID, req.RowID, req.ColumnID, req.Value)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *setCellValueHandler) Path() string {
	return "/tables/:id/set-cell-value"
}

func (h *setCellValueHandler) Method() string {
	return http.MethodPost
}

func (h *setCellValueHandler) AuthRequired() bool {
	return true
}
