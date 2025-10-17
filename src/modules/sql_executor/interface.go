package sql_executor

import (
	"context"
	"database/sql"
)

type ISQLExecutor interface {
	Run(ctx context.Context, dest interface{}, query IToSQL) error
	Exec(ctx context.Context, query IToSQL) (sql.Result, error)
}

type IToSQL interface {
	ToSql() (sqlStr string, args []interface{}, err error)
}
