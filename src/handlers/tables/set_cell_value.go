package tables

import (
	"backend/src/domains/entities"
	"backend/src/handlers"
	"backend/src/modules/web_sockets"
	"backend/src/services"
	"backend/src/services/tables"
	"net/http"

	"github.com/gin-gonic/gin"
)

type setCellValueHandler struct {
	tablesService    services.ITablesService
	databasesService services.IDatabasesService
	tablesHub        *web_sockets.Hub
}

func newSetCellValueHandler(
	tablesHub *web_sockets.Hub,
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
) handlers.IHandler {
	return &setCellValueHandler{
		tablesService:    tablesService,
		databasesService: databasesService,
		tablesHub:        tablesHub,
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

	unlock := h.tablesService.LockTable(tableID)
	defer unlock()
	table, err := h.tablesService.GetTableByID(c, tableID, false)
	if err != nil {
		if tables.IsErrTableNotFound(err) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "table not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(int64)
	authorized, err := h.databasesService.CheckUserRole(c, userID, table.DatabaseID, entities.RoleWriter)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !authorized {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user does not have writer role"})
		return
	}

	var targetColumn *entities.TableColumn
	for _, col := range table.Columns {
		if col.ID == req.ColumnID {
			if col.DeletedAt != nil {
				break
			}
			targetColumn = col
			break
		}
	}
	if targetColumn == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "column not found"})
		return
	}

	if !targetColumn.ValidateColumnValue(req.Value) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid value"})
		return
	}

	err = h.tablesService.SetCellValue(c, userID, table.ID, req.RowID, req.ColumnID, req.Value)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.tablesHub.Broadcast(tableID, entities.EventActionSetCellValue, entities.SetCellValueMessage{
		RowID:    req.RowID,
		ColumnID: req.ColumnID,
		Value:    req.Value,
	})

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
