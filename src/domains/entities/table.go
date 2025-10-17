package entities

import (
	"backend/src/modules/sql_executor"
	"fmt"
	"strings"
	"time"

	"github.com/elgris/sqrl"
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

func (t *Table) AddColumnExpression(column *TableColumn) sql_executor.IToSQL {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("ALTER TABLE %s.%s ADD COLUMN %s text", UsersTablespace, t.ID, column.ID))
	return sqrl.Expr(builder.String())
}

func (t *Table) ReturningCols() []string {
	returningCols := make([]string, 0, 1+len(t.Columns))
	returningCols = append(returningCols, "id")
	for _, col := range t.Columns {
		returningCols = append(returningCols, col.ID)
	}

	return returningCols
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
	Name string     `json:"name"`
	Type ColumnType `json:"type"`
	ID   string     `json:"id"`
}

type ColumnType string

const (
	ColumnTypeText  ColumnType = "text"
	ColumnTypeInt   ColumnType = "int"
	ColumnTypeFloat ColumnType = "float"
	ColumnTypeBool  ColumnType = "bool"
)

type TableRow map[string]any
