package events

import (
	"backend/src/domains/entities"
	"backend/src/handlers"
	"backend/src/modules/web_sockets"
	"backend/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type cellBusyHandler struct {
	tablesHub    *web_sockets.Hub
	usersService services.IUsersService
}

func newCellBusyHandler(
	tablesHub *web_sockets.Hub,
	usersService services.IUsersService,
) handlers.IHandler {
	return &cellBusyHandler{
		tablesHub:    tablesHub,
		usersService: usersService,
	}
}

func (h *cellBusyHandler) Handle(c *gin.Context) {
	req := cellEventRequestDto{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	userID := c.MustGet("user_id").(int64)
	user, err := h.usersService.FindUserByID(c, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	h.tablesHub.Broadcast(req.TableID, entities.EventActionSetCellBusy, &entities.SetCellBusyMessage{
		RowID:    req.RowID,
		ColumnID: req.ColumnID,
		User:     user,
	})

	c.Status(http.StatusOK)
}

func (h *cellBusyHandler) Path() string {
	return "/events/set-cell-busy"
}

func (h *cellBusyHandler) Method() string {
	return http.MethodPost
}

func (h *cellBusyHandler) AuthRequired() bool {
	return true
}
