package databases

import (
	"backend/src/handlers"
	"backend/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createDatabaseHandler struct {
	databasesService services.IDatabasesService
}

func newCreateDatabaseHandler(
	databasesService services.IDatabasesService,
) handlers.IHandler {
	return &createDatabaseHandler{
		databasesService: databasesService,
	}
}

func (h *createDatabaseHandler) Handle(c *gin.Context) {
	req := createDatabaseRequestDto{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	database, err := h.databasesService.AddDatabase(c, c.MustGet("user_id").(int64), req.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newDatabaseResponse(database))
}

func (h *createDatabaseHandler) Path() string {
	return "/databases/create"
}

func (h *createDatabaseHandler) Method() string {
	return http.MethodPost
}

func (h *createDatabaseHandler) AuthRequired() bool {
	return true
}
