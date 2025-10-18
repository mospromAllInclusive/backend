package databases

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

type deleteUserHandler struct {
	databasesService services.IDatabasesService
	tablesService    services.ITablesService
	usersHub         *web_sockets.Hub
}

func newDeleteUserHandler(
	usersHub *web_sockets.Hub,
	databasesService services.IDatabasesService,
	tablesService services.ITablesService,
) handlers.IHandler {
	return &deleteUserHandler{
		databasesService: databasesService,
		tablesService:    tablesService,
		usersHub:         usersHub,
	}
}

func (h *deleteUserHandler) Handle(c *gin.Context) {
	req := deleteUserRequestDto{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	dbID := c.Param("id")
	dbIDInt, err := strconv.ParseInt(dbID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID: " + err.Error()})
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

	err = h.databasesService.DeleteUsersDatabaseRelation(c, req.UserID, dbIDInt)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_ = common.ThrowUserFromDBTables(c, h.tablesService, h.usersHub, req.UserID, dbIDInt)
	h.usersHub.Broadcast(strconv.FormatInt(req.UserID, 10), entities.EventActionFetchDatabases, nil)

	c.Status(http.StatusOK)
}

func (h *deleteUserHandler) Path() string {
	return "/databases/:id/delete-user"
}

func (h *deleteUserHandler) Method() string {
	return http.MethodPost
}

func (h *deleteUserHandler) AuthRequired() bool {
	return true
}
