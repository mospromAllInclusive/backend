package changelog

type cellRequestDto struct {
	TableID  string `json:"table_id" binding:"required"`
	RowID    int64  `json:"row_id" binding:"required"`
	ColumnID string `json:"column_id" binding:"required"`
}
