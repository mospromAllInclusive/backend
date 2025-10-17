package tables

import (
	"backend/src/domains/entities"
	"backend/src/handlers"
	"backend/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addRowHandler struct {
	tablesService    services.ITablesService
	databasesService services.IDatabasesService
}

func newAddRowHandler(
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
) handlers.IHandler {
	return &addRowHandler{
		tablesService:    tablesService,
		databasesService: databasesService,
	}
}

func (h *addRowHandler) Handle(c *gin.Context) {
	req := addRowRequestDto{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	tableID := c.Param("id")
	if tableID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid table id"})
		return
	}

	table, err := h.tablesService.GetTableByID(c, tableID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if table == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "table not found"})
		return
	}

	authorized, err := h.databasesService.CheckUserRole(c, c.MustGet("user_id").(int64), table.DatabaseID, entities.RoleWriter)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !authorized {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user does not have writer role"})
		return
	}

	row, err := h.tablesService.AddRow(c, table, req.SortIndex)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newRowResponse(row))
}

func (h *addRowHandler) Path() string {
	return "/tables/:id/add-row"
}

func (h *addRowHandler) Method() string {
	return http.MethodPost
}

func (h *addRowHandler) AuthRequired() bool {
	return true
}
