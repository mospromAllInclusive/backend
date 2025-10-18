package tables

import (
	"backend/src/domains/entities"
	"backend/src/handlers"
	"backend/src/services"
	"backend/src/services/tables"
	"net/http"

	"github.com/gin-gonic/gin"
)

const filename = "export.xlsx"

type exportTableHandler struct {
	tablesService    services.ITablesService
	databasesService services.IDatabasesService
}

func newExportTableHandler(
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
) handlers.IHandler {
	return &exportTableHandler{
		tablesService:    tablesService,
		databasesService: databasesService,
	}
}

func (h *exportTableHandler) Handle(c *gin.Context) {
	tableID := c.Param("id")
	if tableID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid table id"})
		return
	}

	table, err := h.tablesService.GetTableByID(c, tableID, false)
	if err != nil {
		if tables.IsErrTableNotFound(err) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "table not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	authorized, err := h.databasesService.CheckUserRole(c, c.MustGet("user_id").(int64), table.DatabaseID, entities.RoleReader)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !authorized {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user does not have reader role"})
		return
	}

	file, err := h.tablesService.ExportTable(c, table)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", `attachment; filename="`+filename+`"; filename*=UTF-8''`+filename)
	c.Header("X-Content-Type-Options", "nosniff")
	c.Status(http.StatusOK)

	if _, err := file.WriteTo(c.Writer); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (h *exportTableHandler) Path() string {
	return "/tables/:id/export"
}

func (h *exportTableHandler) Method() string {
	return http.MethodGet
}

func (h *exportTableHandler) AuthRequired() bool {
	return true
}
