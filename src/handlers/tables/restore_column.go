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

type restoreColumnHandler struct {
	tablesService    services.ITablesService
	databasesService services.IDatabasesService
	tablesHub        *web_sockets.Hub
}

func newRestoreColumnHandler(
	tablesHub *web_sockets.Hub,
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
) handlers.IHandler {
	return &restoreColumnHandler{
		tablesService:    tablesService,
		databasesService: databasesService,
		tablesHub:        tablesHub,
	}
}

func (h *restoreColumnHandler) Handle(c *gin.Context) {
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

	authorized, err := h.databasesService.CheckUserRole(c, c.MustGet("user_id").(int64), table.DatabaseID, entities.RoleAdmin)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !authorized {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user does not have admin role"})
		return
	}

	table, err = h.tablesService.RestoreColumn(c, req.ColumnID, req.TableID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.tablesHub.Broadcast(req.TableID, entities.EventActionFetchTable, nil)
	c.JSON(http.StatusOK, common.NewTableResponse(table))
}

func (h *restoreColumnHandler) Path() string {
	return "/tables/restore-column"
}

func (h *restoreColumnHandler) Method() string {
	return http.MethodPost
}

func (h *restoreColumnHandler) AuthRequired() bool {
	return true
}
