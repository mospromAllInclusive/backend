package databases

import (
	"backend/src/handlers"
	"backend/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type listDatabasesHandler struct {
	databasesService services.IDatabasesService
	tablesService    services.ITablesService
}

func newListDatabasesHandler(
	databasesService services.IDatabasesService,
	tablesService services.ITablesService,
) handlers.IHandler {
	return &listDatabasesHandler{
		databasesService: databasesService,
		tablesService:    tablesService,
	}
}

func (h *listDatabasesHandler) Handle(c *gin.Context) {
	usersDatabases, err := h.databasesService.GetUsersDatabases(c, c.MustGet("user_id").(int64))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dbIDs := make([]int64, 0, len(usersDatabases))
	for _, db := range usersDatabases {
		dbIDs = append(dbIDs, db.DatabaseID)
	}

	tables, err := h.tablesService.ListByDatabaseIDs(c, dbIDs)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newDatabaseListResponse(usersDatabases, tables))
}

func (h *listDatabasesHandler) Path() string {
	return "/databases/list"
}

func (h *listDatabasesHandler) Method() string {
	return http.MethodGet
}

func (h *listDatabasesHandler) AuthRequired() bool {
	return true
}
