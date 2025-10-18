package databases

import (
	"backend/src/domains/entities"
	"backend/src/handlers"
	"backend/src/modules/web_sockets"
	"backend/src/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type createDatabaseHandler struct {
	usersHub         *web_sockets.Hub
	databasesService services.IDatabasesService
}

func newCreateDatabaseHandler(
	usersHub *web_sockets.Hub,
	databasesService services.IDatabasesService,
) handlers.IHandler {
	return &createDatabaseHandler{
		usersHub:         usersHub,
		databasesService: databasesService,
	}
}

func (h *createDatabaseHandler) Handle(c *gin.Context) {
	req := createDatabaseRequestDto{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	userID := c.MustGet("user_id").(int64)
	database, err := h.databasesService.AddDatabase(c, userID, req.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.usersHub.Broadcast(strconv.FormatInt(userID, 10), entities.EventActionFetchDatabases, nil)

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
