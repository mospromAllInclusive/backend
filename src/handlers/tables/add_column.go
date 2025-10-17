package tables

import (
	"backend/src/domains/entities"
	"backend/src/handlers"
	"backend/src/handlers/common"
	"backend/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addColumnHandler struct {
	tablesService    services.ITablesService
	databasesService services.IDatabasesService
}

func newAddColumnHandler(
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
) handlers.IHandler {
	return &addColumnHandler{
		tablesService:    tablesService,
		databasesService: databasesService,
	}
}

func (h *addColumnHandler) Handle(c *gin.Context) {
	req := addColumnRequestDto{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	table, err := h.tablesService.GetTableByID(c, req.TableID)
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

	col := req.Column.toEntity()
	table, err = h.tablesService.AddColumnToTable(c, col, req.TableID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
