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

type editColumnHandler struct {
	tablesService    services.ITablesService
	databasesService services.IDatabasesService
	hub              *web_sockets.Hub
}

func newEditColumnHandler(
	hub *web_sockets.Hub,
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
) handlers.IHandler {
	return &editColumnHandler{
		tablesService:    tablesService,
		databasesService: databasesService,
		hub:              hub,
	}
}

func (h *editColumnHandler) Handle(c *gin.Context) {
	req := editColumnRequestDto{}
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

	col, err := req.Column.toEntity()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	columnExists := false
	for _, tableColumn := range table.Columns {
		if tableColumn.ID == col.ID {
			if tableColumn.DeletedAt != nil {
				break
			}
			columnExists = true
			break
		}
	}

	if !columnExists {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "column not found"})
		return
	}

	invalidValues, err := h.tablesService.ValidateColumnValues(c, req.TableID, col)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(invalidValues) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, invalidColumValuesResponse{InvalidValues: invalidValues})
		return
	}

	table, _, err = h.tablesService.EditTableColumn(c, col, req.TableID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.hub.Broadcast(req.TableID, entities.EventActionFetchTable, nil)
	c.JSON(http.StatusOK, common.NewTableResponse(table))
}

func (h *editColumnHandler) Path() string {
	return "/tables/edit-column"
}

func (h *editColumnHandler) Method() string {
	return http.MethodPost
}

func (h *editColumnHandler) AuthRequired() bool {
	return true
}
