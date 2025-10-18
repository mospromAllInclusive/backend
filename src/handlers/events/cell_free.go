package events

import (
	"backend/src/domains/entities"
	"backend/src/handlers"
	"backend/src/modules/web_sockets"
	"net/http"

	"github.com/gin-gonic/gin"
)

type cellFreeHandler struct {
	tablesHub *web_sockets.Hub
}

func newCellFreeHandler(
	tablesHub *web_sockets.Hub,
) handlers.IHandler {
	return &cellFreeHandler{
		tablesHub: tablesHub,
	}
}

func (h *cellFreeHandler) Handle(c *gin.Context) {
	req := cellEventRequestDto{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	userID := c.MustGet("user_id").(int64)
	h.tablesHub.Broadcast(req.TableID, entities.EventActionSetCellFree, &entities.SetCellFreeMessage{
		RowID:    req.RowID,
		ColumnID: req.ColumnID,
		UserID:   userID,
	})

	c.Status(http.StatusOK)
}

func (h *cellFreeHandler) Path() string {
	return "/events/set-cell-free"
}

func (h *cellFreeHandler) Method() string {
	return http.MethodPost
}

func (h *cellFreeHandler) AuthRequired() bool {
	return true
}
