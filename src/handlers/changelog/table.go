package changelog

import (
	"backend/src/domains/entities"
	"backend/src/handlers"
	"backend/src/services"
	"backend/src/services/tables"
	"net/http"

	"github.com/gin-gonic/gin"
)

type tableHandler struct {
	changelogService services.IChangelogService
	tablesService    services.ITablesService
	databasesService services.IDatabasesService
}

func newTableHandler(
	changelogService services.IChangelogService,
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
) handlers.IHandler {
	return &tableHandler{
		changelogService: changelogService,
		tablesService:    tablesService,
		databasesService: databasesService,
	}
}

func (h *tableHandler) Handle(c *gin.Context) {
	req := tableRequestDto{}
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

	changelog, err := h.changelogService.ListChangelogForTable(c, req.TableID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newTableChangelogResponse(changelog))
}

func (h *tableHandler) Path() string {
	return "/changelog/table"
}

func (h *tableHandler) Method() string {
	return http.MethodPost
}

func (h *tableHandler) AuthRequired() bool {
	return true
}
