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

type addColumnHandler struct {
	hub              *web_sockets.Hub
	tablesService    services.ITablesService
	databasesService services.IDatabasesService
}

func newAddColumnHandler(
	hub *web_sockets.Hub,
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
) handlers.IHandler {
	return &addColumnHandler{
		tablesService:    tablesService,
		databasesService: databasesService,
		hub:              hub,
	}
}

func (h *addColumnHandler) Handle(c *gin.Context) {
	req := addColumnRequestDto{}
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

	col, err := req.Column.toEntity()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	table, err = h.tablesService.AddColumnToTable(c, col, req.TableID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.hub.Broadcast(req.TableID, entities.EventActionFetchTable, nil)
	c.JSON(http.StatusOK, common.NewTableResponse(table))
}

func (h *addColumnHandler) Path() string {
	return "/tables/add-column"
}

func (h *addColumnHandler) Method() string {
	return http.MethodPost
}

func (h *addColumnHandler) AuthRequired() bool {
	return true
}
