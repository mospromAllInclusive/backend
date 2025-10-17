package tables

import "backend/src/domains/entities"

type createTableRequestDto struct {
	Name       string   `json:"name" binding:"required"`
	DatabaseID int64    `json:"database_id" binding:"required"`
	Columns    []column `json:"columns" binding:"required"`
}

func (c *createTableRequestDto) toEntity() *entities.Table {
	cols := make([]*entities.TableColumn, 0, len(c.Columns))
	for _, col := range c.Columns {
		cols = append(cols, col.toEntity())
	}
	return &entities.Table{
		Name:       c.Name,
		DatabaseID: c.DatabaseID,
		Columns:    cols,
	}
}

type column struct {
	Name string              `json:"name" binding:"required"`
	Type entities.ColumnType `json:"type" binding:"required"`
}

func (c *column) toEntity() *entities.TableColumn {
	return &entities.TableColumn{
		Name: c.Name,
		Type: c.Type,
	}
}

type requestByTableID struct {
	TableID string `json:"table_id" binding:"required"`
}

type addColumnRequestDto struct {
	TableID string `json:"table_id" binding:"required"`
	Column  column `json:"column" binding:"required"`
}

type defaultColumnRequestDto struct {
	TableID  string `json:"table_id" binding:"required"`
	ColumnID string `json:"column_id" binding:"required"`
}

type addRowRequestDto struct {
	SortIndex *int64             `json:"sort_index"`
	Data      map[string]*string `json:"data"`
}

type defaultRowRequestDto struct {
	RowID int64 `json:"row_id" binding:"required"`
}

type moveRowRequestDto struct {
	defaultRowRequestDto
	SortIndex int64 `json:"sort_index" binding:"required"`
}

type setCellValueRequestDto struct {
	defaultRowRequestDto
	ColumnID string  `json:"column_id" binding:"required"`
	Value    *string `json:"value"`
}
