package tables

import (
	"backend/src/domains/entities"
	"backend/src/handlers"
	"backend/src/modules/web_sockets"
	"backend/src/services"
	"backend/src/services/tables"
	"net/http"

	"github.com/gin-gonic/gin"
)

type addRowHandler struct {
	tablesHub        *web_sockets.Hub
	tablesService    services.ITablesService
	databasesService services.IDatabasesService
	changelogService services.IChangelogService
}

func newAddRowHandler(
	tablesHub *web_sockets.Hub,
	tablesService services.ITablesService,
	databasesService services.IDatabasesService,
	changelogService services.IChangelogService,
) handlers.IHandler {
	return &addRowHandler{
		tablesService:    tablesService,
		databasesService: databasesService,
		changelogService: changelogService,
		tablesHub:        tablesHub,
	}
}

func (h *addRowHandler) Handle(c *gin.Context) {
	req := addRowRequestDto{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	tableID := c.Param("id")
	if tableID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid table id"})
		return
	}

	unlock := h.tablesService.LockTable(tableID)
	defer unlock()
	table, err := h.tablesService.GetTableByID(c, tableID, false)
	if err != nil {
		if tables.IsErrTableNotFound(err) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "table not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(int64)
	authorized, err := h.databasesService.CheckUserRole(c, userID, table.DatabaseID, entities.RoleWriter)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !authorized {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user does not have writer role"})
		return
	}

	existentColumns := make(map[string]*entities.TableColumn, len(table.Columns))
	for _, col := range table.Columns {
		if col.DeletedAt != nil {
			continue
		}
		existentColumns[col.ID] = col
	}

	invalidValues := make([]*string, 0)
	for colID, value := range req.Data {
		column, ok := existentColumns[colID]
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "column " + colID + " does not exist"})
			return
		}

		if !column.ValidateColumnValue(value) {
			invalidValues = append(invalidValues, value)
		}
	}

	if len(invalidValues) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, invalidColumValuesResponse{InvalidValues: invalidValues})
		return
	}

	row, err := h.tablesService.AddRow(c, userID, table, req.Data, req.SortIndex)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.tablesHub.Broadcast(tableID, entities.EventActionFetchTable, nil)

	rowChange := &entities.RowChange{
		ChangeType: entities.ChangeTypeAdd,
		Before:     nil,
		After:      entities.NewRowInfoForChangelog(table, row),
	}

	changelogItem := rowChange.ToChangelogItem(userID, table.ID, row.GetID())
	err = h.changelogService.WriteChangelog(c, changelogItem)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newRowResponse(row))
}

func (h *addRowHandler) Path() string {
	return "/tables/:id/add-row"
}

func (h *addRowHandler) Method() string {
	return http.MethodPost
}

func (h *addRowHandler) AuthRequired() bool {
	return true
}
