package tables

import (
	"backend/src/domains/entities"
	"backend/src/handlers"
	"backend/src/handlers/common"
	"backend/src/modules/web_sockets"
	"backend/src/services"
	"backend/src/services/tables"
	"net/http"

	"github.com/gin-gonic/gin"
)

type deleteColumnHandler struct {
	tablesService    services.ITablesService
	databasesService services.IDatabasesService
	changelogService services.IChangelogService
	tablesHub        *web_sockets.Hub
}

func newDeleteColumnHandler(
	tablesHub *web_sockets.Hub,
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
	changelogService services.IChangelogService,
) handlers.IHandler {
	return &deleteColumnHandler{
		tablesService:    tablesService,
		databasesService: databasesService,
		changelogService: changelogService,
		tablesHub:        tablesHub,
	}
}

func (h *deleteColumnHandler) Handle(c *gin.Context) {
	req := defaultColumnRequestDto{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	unlock := h.tablesService.LockTable(req.TableID)
	defer unlock()
	table, err := h.tablesService.GetTableByID(c, req.TableID, false)
	if err != nil {
		if tables.IsErrTableNotFound(err) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "table not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(int64)
	authorized, err := h.databasesService.CheckUserRole(c, userID, table.DatabaseID, entities.RoleAdmin)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !authorized {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user does not have admin role"})
		return
	}

	table, err = h.tablesService.DeleteColumn(c, req.ColumnID, req.TableID)
	if err != nil {
		if tables.IsErrColumnNotFound(err) {
			c.JSON(http.StatusOK, common.NewTableResponse(table))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.tablesHub.Broadcast(req.TableID, entities.EventActionFetchTable, nil)

	var deletedCol *entities.TableColumn
	for _, col := range table.Columns {
		if col.ID == req.ColumnID {
			deletedCol = col
			break
		}
	}

	columnChange := &entities.ColumnChange{
		ChangeType: entities.ChangeTypeDelete,
		Before:     deletedCol,
		After:      nil,
	}

	changelogItem := columnChange.ToChangelogItem(userID, table.ID, deletedCol.ID)
	err = h.changelogService.WriteChangelog(c, changelogItem)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, common.NewTableResponse(table))
}

func (h *deleteColumnHandler) Path() string {
	return "/tables/delete-column"
}

func (h *deleteColumnHandler) Method() string {
	return http.MethodPost
}

func (h *deleteColumnHandler) AuthRequired() bool {
	return true
}
