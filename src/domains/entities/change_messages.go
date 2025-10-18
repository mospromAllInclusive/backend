package entities

const (
	EventActionSetCellValue    string = "set_cell_value"
	EventActionFetchTable      string = "fetch_table"
	EventActionGoAwayFromTable string = "go_away_from_table"
	EventActionFetchDatabases  string = "fetch_databases"
	EventActionSetCellBusy     string = "set_cell_busy"
	EventActionSetCellFree     string = "set_cell_free"
)

type SetCellValueMessage struct {
	RowID    int64   `json:"row_id"`
	ColumnID string  `json:"column_id"`
	Value    *string `json:"value"`
}

type GoAwayFromTableMessage struct {
	TableID string `json:"table_id"`
}

type SetCellBusyMessage struct {
	RowID    int64  `json:"row_id"`
	ColumnID string `json:"column_id"`
	User     *User  `json:"user"`
}

type SetCellFreeMessage struct {
	RowID    int64  `json:"row_id"`
	ColumnID string `json:"column_id"`
	UserID   int64  `json:"user_id"`
}
