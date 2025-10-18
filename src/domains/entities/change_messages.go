package entities

const (
	EventActionSetCellValue string = "set_cell_value"
	EventActionFetchTable   string = "fetch_table"
)

type SetCellValueMessage struct {
	RowID    int64   `json:"row_id"`
	ColumnID string  `json:"column_id"`
	Value    *string `json:"value"`
}
