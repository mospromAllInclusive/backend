package tables

import (
	"backend/src/domains/entities"
	"backend/src/handlers"
	"backend/src/services"
	"backend/src/services/tables"
	"net/http"

	"github.com/gin-gonic/gin"
)

type deleteTableHandler struct {
	tablesService    services.ITablesService
	databasesService services.IDatabasesService
}

func newDeleteTableHandler(
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
) handlers.IHandler {
	return &deleteTableHandler{
		tablesService:    tablesService,
		databasesService: databasesService,
	}
}

func (h *deleteTableHandler) Handle(c *gin.Context) {
	req := requestByTableID{}
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

	err = h.tablesService.DeleteTable(c, req.TableID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *deleteTableHandler) Path() string {
	return "/tables/delete"
}

func (h *deleteTableHandler) Method() string {
	return http.MethodPost
}

func (h *deleteTableHandler) AuthRequired() bool {
	return true
}
