package changelog

import (
	"backend/src/domains/entities"
	"backend/src/handlers"
	"backend/src/services"
	"backend/src/services/tables"
	"net/http"

	"github.com/gin-gonic/gin"
)

type cellHandler struct {
	changelogService services.IChangelogService
	tablesService    services.ITablesService
	databasesService services.IDatabasesService
}

func newCellHandler(
	changelogService services.IChangelogService,
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
) handlers.IHandler {
	return &cellHandler{
		changelogService: changelogService,
		tablesService:    tablesService,
		databasesService: databasesService,
	}
}

func (h *cellHandler) Handle(c *gin.Context) {
	req := cellRequestDto{}
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

	authorized, err := h.databasesService.CheckUserRole(c, c.MustGet("user_id").(int64), table.DatabaseID, entities.RoleReader)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !authorized {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user does not have reader role"})
		return
	}

	changelog, err := h.changelogService.ListChangelogForCell(c, req.TableID, req.ColumnID, req.RowID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newCellChangelogResponse(changelog))
}

func (h *cellHandler) Path() string {
	return "/changelog/cell"
}

func (h *cellHandler) Method() string {
	return http.MethodPost
}

func (h *cellHandler) AuthRequired() bool {
	return true
}
