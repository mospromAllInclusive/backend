package entities

import (
	"backend/src/modules/sql_executor"
	"fmt"
	"strings"
	"time"

	"github.com/elgris/sqrl"
	"github.com/elgris/sqrl/pg"
)

const (
	UsersTablespace = "users_tablespace"
)

type Table struct {
	ID         string
	Name       string
	DatabaseID int64
	Columns    []*TableColumn
	CreatedAt  time.Time
}

func (t *Table) ToDBTable() *DBTable {
	columns := make([]*TableColumn, 0)
	if t.Columns != nil {
		columns = t.Columns
	}
	return &DBTable{
		ID:         t.ID,
		Name:       t.Name,
		DatabaseID: t.DatabaseID,
		Columns:    JSONB[[]*TableColumn]{v: &columns},
		CreatedAt:  t.CreatedAt,
	}
}

func (t *Table) CreateExpression() sql_executor.IToSQL {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("CREATE TABLE %s.%s ( ", UsersTablespace, t.ID))
	builder.WriteString("id bigserial primary key, sort_index bigserial not null, sort_index_version bigint not null default 0, ")

	for _, column := range t.Columns {
		builder.WriteString(fmt.Sprintf("%s text, ", column.ID))
	}
	builder.WriteString("deleted_at timestamp with time zone)")

	return sqrl.Expr(builder.String())
}

func (t *Table) CreateSortIndexExpression() sql_executor.IToSQL {
	return sqrl.Expr(fmt.Sprintf("create index if not exists %s_order_idx on %s.%s (sort_index asc, sort_index_version desc)", t.ID, UsersTablespace, t.ID))
}

func (t *Table) CreateColumnIndexExpression(columnID string) sql_executor.IToSQL {
	return sqrl.Expr(fmt.Sprintf("create index if not exists %s_%s_trgm_idx on %s.%s USING gin (%s gin_trgm_ops)", t.ID, columnID, UsersTablespace, t.ID, columnID))
}

func (t *Table) AddColumnExpression(column *TableColumn) sql_executor.IToSQL {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("ALTER TABLE %s.%s ADD COLUMN %s text", UsersTablespace, t.ID, column.ID))
	return sqrl.Expr(builder.String())
}

func (t *Table) ReturningCols() []string {
	returningCols := make([]string, 0, 1+len(t.Columns))
	returningCols = append(returningCols, "id")
	for _, col := range t.Columns {
		if col.DeletedAt != nil {
			continue
		}
		returningCols = append(returningCols, col.ID)
	}

	return returningCols
}

func (t *Table) ValidateParams(params *ReadTableParams) bool {
	return t.ValidateFilter(params) && t.ValidateSort(params)
}

func (t *Table) ValidateFilter(params *ReadTableParams) bool {
	if params.FilterBy == nil {
		return true
	}

	for _, col := range t.Columns {
		if col.DeletedAt != nil {
			continue
		}
		if col.ID == *params.FilterBy {
			return true
		}
	}

	return false
}

func (t *Table) ValidateSort(params *ReadTableParams) bool {
	if params.SortBy == nil {
		return true
	}

	for _, col := range t.Columns {
		if col.DeletedAt != nil {
			continue
		}
		if col.ID == *params.SortBy {
			return true
		}
	}

	return false
}

type DBTable struct {
	Name       string                `db:"name"`
	ID         string                `db:"id"`
	DatabaseID int64                 `db:"database_id"`
	Columns    JSONB[[]*TableColumn] `db:"columns"`
	CreatedAt  time.Time             `db:"created_at"`
}

func (t *DBTable) ToTable() *Table {
	return &Table{
		Name:       t.Name,
		ID:         t.ID,
		DatabaseID: t.DatabaseID,
		Columns:    *t.Columns.Get(),
		CreatedAt:  t.CreatedAt,
	}
}

type TableColumn struct {
	Name      string     `json:"name"`
	Type      ColumnType `json:"type"`
	Enum      []string   `json:"enum,omitempty"`
	ID        string     `json:"id"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (c *TableColumn) NeedToBeUpdated(new *TableColumn) bool {
	if c.Type != new.Type || c.Name != new.Name {
		return true
	}

	newEnum := make(map[string]struct{})
	for _, v := range new.Enum {
		newEnum[v] = struct{}{}
	}

	if len(c.Enum) != len(newEnum) {
		return true
	}

	for _, k := range c.Enum {
		if _, ok := newEnum[k]; !ok {
			return true
		}
	}

	return false
}

type ColumnType string

const (
	ColumnTypeText      ColumnType = "text"
	ColumnTypeNumeric   ColumnType = "numeric"
	ColumnTypeEnum      ColumnType = "enum"
	ColumnTypeTimestamp ColumnType = "timestamp"
)

func (t ColumnType) TypeCast() string {
	switch t {
	case ColumnTypeNumeric:
		return "::numeric"
	default:
		return ""
	}
}

type TableRow map[string]any

func (t TableRow) GetID() int64 {
	return t["id"].(int64)
}

type ReadTableParams struct {
	Page        int     `form:"page" binding:"required,min=1"`
	PerPage     int     `form:"perPage" binding:"required,min=1,max=1000"`
	SortBy      *string `form:"sortBy" binding:"omitempty,gt=0"`
	SortDir     *string `form:"sortDir" binding:"omitempty,oneof=asc desc"`
	FilterBy    *string `form:"filterBy" binding:"omitempty,gt=0,excluded_without=FilterValue"`
	FilterValue *string `form:"filterValue" binding:"omitempty,gt=0"`
}

func (p ReadTableParams) GetLimit() int {
	return p.PerPage
}

func (p ReadTableParams) GetOffset() int {
	return (p.Page - 1) * p.PerPage
}

func (p ReadTableParams) GetSortBy(t *Table) string {
	if p.SortBy == nil {
		return ""
	}

	sortDir := "asc"
	if p.SortDir != nil && *p.SortDir == "desc" {
		sortDir = "desc"
	}

	var columnType ColumnType
	for _, col := range t.Columns {
		if col.ID == *p.SortBy {
			columnType = col.Type
			break
		}
	}

	return fmt.Sprintf("%s%s %s", *p.SortBy, columnType.TypeCast(), sortDir)
}

func (p ReadTableParams) GetFilter() (bool, string, interface{}) {
	if p.FilterBy == nil {
		return false, "", nil
	}

	filterSql, filterValue := LikeFilter(*p.FilterValue, *p.FilterBy)
	return true, filterSql, filterValue
}

func escapeLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

func LikeFilter(value string, field string) (string, interface{}) {
	escaped := fmt.Sprintf("%%%s%%", escapeLike(value))
	filterSql := fmt.Sprintf("%s ILIKE ANY (?)", field)
	return filterSql, pg.Array([]string{escaped})
}
