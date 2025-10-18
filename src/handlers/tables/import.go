package tables

import (
	"backend/src/domains/entities"
	"backend/src/handlers"
	"backend/src/handlers/common"
	"backend/src/modules/web_sockets"
	"backend/src/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type importTableHandler struct {
	hub              *web_sockets.Hub
	tablesService    services.ITablesService
	databasesService services.IDatabasesService
	fileService      services.IFileService
}

func newImportTableHandler(
	hub *web_sockets.Hub,
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
	fileService services.IFileService,
) handlers.IHandler {
	return &importTableHandler{
		hub:              hub,
		tablesService:    tablesService,
		databasesService: databasesService,
		fileService:      fileService,
	}
}

func (h *importTableHandler) Handle(c *gin.Context) {
	dbID := c.PostForm("database_id")
	if dbID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing database_id"})
		return
	}

	dbIDInt, err := strconv.ParseInt(dbID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid database_id: " + err.Error()})
		return
	}

	authorized, err := h.databasesService.CheckUserRole(c, c.MustGet("user_id").(int64), dbIDInt, entities.RoleAdmin)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !authorized {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user does not have admin role"})
		return
	}

	tableName := c.PostForm("table_name")
	if tableName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing table_name"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing file"})
		return
	}

	columns, data, err := h.fileService.ReadFile(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	table, err := h.tablesService.ImportTable(c, tableName, dbIDInt, columns, data)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, common.NewTableResponse(table))
}

func (h *importTableHandler) Path() string {
	return "/tables/import"
}

func (h *importTableHandler) Method() string {
	return http.MethodPost
}

func (h *importTableHandler) AuthRequired() bool {
	return true
}
