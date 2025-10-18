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

type restoreTableHandler struct {
	tablesService    services.ITablesService
	databasesService services.IDatabasesService
	usersHub         *web_sockets.Hub
}

func newRestoreTableHandler(
	usersHub *web_sockets.Hub,
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
) handlers.IHandler {
	return &restoreTableHandler{
		tablesService:    tablesService,
		databasesService: databasesService,
		usersHub:         usersHub,
	}
}

func (h *restoreTableHandler) Handle(c *gin.Context) {
	req := requestByTableID{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	table, err := h.tablesService.GetTableByID(c, req.TableID, true)
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

	unlock := h.tablesService.LockTable(req.TableID)
	defer unlock()
	err = h.tablesService.RestoreTable(c, req.TableID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_ = common.SendActionToDBUsers(c, h.databasesService, h.usersHub, table.DatabaseID, entities.EventActionFetchDatabases)

	c.Status(http.StatusOK)
}

func (h *restoreTableHandler) Path() string {
	return "/tables/restore"
}

func (h *restoreTableHandler) Method() string {
	return http.MethodPost
}

func (h *restoreTableHandler) AuthRequired() bool {
	return true
}
