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

type deleteRowHandler struct {
	tablesService    services.ITablesService
	databasesService services.IDatabasesService
	changelogService services.IChangelogService
	hub              *web_sockets.Hub
}

func newDeleteRowHandler(
	hub *web_sockets.Hub,
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
	changelogService services.IChangelogService,
) handlers.IHandler {
	return &deleteRowHandler{
		tablesService:    tablesService,
		databasesService: databasesService,
		changelogService: changelogService,
		hub:              hub,
	}
}

func (h *deleteRowHandler) Handle(c *gin.Context) {
	req := defaultRowRequestDto{}
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

	row, err := h.tablesService.DeleteRow(c, table.ID, req.RowID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if row == nil {
		c.Status(http.StatusOK)
		return
	}

	h.hub.Broadcast(tableID, entities.EventActionFetchTable, nil)

	rowChange := &entities.RowChange{
		ChangeType: entities.ChangeTypeDelete,
		Before:     entities.NewRowInfoForChangelog(table, row),
		After:      nil,
	}

	changelogItem := rowChange.ToChangelogItem(userID, table.ID, row.GetID())
	err = h.changelogService.WriteChangelog(c, changelogItem)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *deleteRowHandler) Path() string {
	return "/tables/:id/delete-row"
}

func (h *deleteRowHandler) Method() string {
	return http.MethodPost
}

func (h *deleteRowHandler) AuthRequired() bool {
	return true
}
