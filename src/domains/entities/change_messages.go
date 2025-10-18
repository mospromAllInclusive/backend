package entities

const (
	EventActionSetCellValue    string = "set_cell_value"
	EventActionFetchTable      string = "fetch_table"
	EventActionGoAwayFromTable string = "go_away_from_table"
	EventActionFetchDatabases  string = "fetch_databases"
)

type SetCellValueMessage struct {
	RowID    int64   `json:"row_id"`
	ColumnID string  `json:"column_id"`
	Value    *string `json:"value"`
}

type GoAwayFromTableMessage struct {
	TableID string `json:"table_id"`
}
