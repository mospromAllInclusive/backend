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
	hub              *web_sockets.Hub
}

func newDeleteColumnHandler(
	hub *web_sockets.Hub,
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
) handlers.IHandler {
	return &deleteColumnHandler{
		tablesService:    tablesService,
		databasesService: databasesService,
		hub:              hub,
	}
}

func (h *deleteColumnHandler) Handle(c *gin.Context) {
	req := defaultColumnRequestDto{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

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

	table, err = h.tablesService.DeleteColumn(c, req.ColumnID, req.TableID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.hub.Broadcast(req.TableID, entities.EventActionFetchTable, nil)
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
