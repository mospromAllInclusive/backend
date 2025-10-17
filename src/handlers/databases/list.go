package databases

import (
	"backend/src/handlers"
	"backend/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type listDatabasesHandler struct {
	databasesService services.IDatabasesService
}

func newListDatabasesHandler(
	databasesService services.IDatabasesService,
) handlers.IHandler {
	return &listDatabasesHandler{
		databasesService: databasesService,
	}
}

func (h *listDatabasesHandler) Handle(c *gin.Context) {
	usersDatabases, err := h.databasesService.GetUsersDatabases(c, c.MustGet("user_id").(int64))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newDatabaseListResponse(usersDatabases))
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
