package tables

import (
	"backend/src/domains/entities"
	"backend/src/handlers"
	"backend/src/handlers/common"
	"backend/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createTableHandler struct {
	tablesService    services.ITablesService
	databasesService services.IDatabasesService
}

func newCreateTableHandler(
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
) handlers.IHandler {
	return &createTableHandler{
		tablesService:    tablesService,
		databasesService: databasesService,
	}
}

func (h *createTableHandler) Handle(c *gin.Context) {
	req := createTableRequestDto{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	authorized, err := h.databasesService.CheckUserRole(c, c.MustGet("user_id").(int64), req.DatabaseID, entities.RoleAdmin)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !authorized {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user does not have admin role"})
		return
	}

	table := req.toEntity()
	table, err = h.tablesService.CreateTable(c, table)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, common.NewTableResponse(table))
}

func (h *createTableHandler) Path() string {
	return "/tables/create"
}

func (h *createTableHandler) Method() string {
	return http.MethodPost
}

func (h *createTableHandler) AuthRequired() bool {
	return true
}
