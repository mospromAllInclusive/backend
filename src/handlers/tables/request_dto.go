package tables

import (
	"backend/src/domains/entities"
	"fmt"
)

type createTableRequestDto struct {
	Name       string   `json:"name" binding:"required"`
	DatabaseID int64    `json:"database_id" binding:"required"`
	Columns    []column `json:"columns" binding:"required"`
}

func (c *createTableRequestDto) toEntity() (*entities.Table, error) {
	cols := make([]*entities.TableColumn, 0, len(c.Columns))
	for _, col := range c.Columns {
		entity, err := col.toEntity()
		if err != nil {
			return nil, err
		}
		cols = append(cols, entity)
	}
	return &entities.Table{
		Name:       c.Name,
		DatabaseID: c.DatabaseID,
		Columns:    cols,
	}, nil
}

type column struct {
	Name string              `json:"name" binding:"required"`
	Type entities.ColumnType `json:"type" binding:"required,oneof=text numeric enum timestamp"`
	Enum []string            `json:"enum" binding:"omitempty,dive,required"`
}

func (c *column) DistinctEnum() {
	if c.Type != entities.ColumnTypeEnum {
		c.Enum = nil
	}

	seen := make(map[string]bool)
	enum := make([]string, 0, len(c.Enum))
	for _, v := range c.Enum {
		if seen[v] {
			continue
		}
		seen[v] = true
		enum = append(enum, v)
	}

	c.Enum = enum
}

func (c *column) toEntity() (*entities.TableColumn, error) {
	if c.Type == entities.ColumnTypeEnum && len(c.Enum) == 0 {
		return nil, fmt.Errorf("enum column must have at least one value")
	}
	c.DistinctEnum()
	return &entities.TableColumn{
		Name: c.Name,
		Type: c.Type,
		Enum: c.Enum,
	}, nil
}

type columnWithID struct {
	ID string `json:"id" binding:"required"`
	column
}

func (c *columnWithID) toEntity() (*entities.TableColumn, error) {
	if c.Type == entities.ColumnTypeEnum && len(c.Enum) == 0 {
		return nil, fmt.Errorf("enum column must have at least one value")
	}
	c.DistinctEnum()
	return &entities.TableColumn{
		ID:   c.ID,
		Name: c.Name,
		Type: c.Type,
		Enum: c.Enum,
	}, nil
}

type requestByTableID struct {
	TableID string `json:"table_id" binding:"required"`
}

type addColumnRequestDto struct {
	TableID string `json:"table_id" binding:"required"`
	Column  column `json:"column" binding:"required"`
}

type editColumnRequestDto struct {
	TableID string       `json:"table_id" binding:"required"`
	Column  columnWithID `json:"column" binding:"required"`
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
