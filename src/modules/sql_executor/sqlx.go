package sql_executor

import (
	"context"
	"database/sql"
	"errors"
	"github.com/georgysavva/scany/v2/dbscan"
	"github.com/jmoiron/sqlx"
	"log"
	"reflect"

	_ "github.com/lib/pq"
)

type sqlxExecutor struct {
	db      *sqlx.DB
	scanner *dbscan.API
}

func NewSQlXExecutor(driverName string, dataSourceName string) ISQLExecutor {
	dbConn, err := sqlx.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	dbConn.SetMaxIdleConns(10)
	dbConn.SetMaxOpenConns(10)

	scanner, err := dbscan.NewAPI(dbscan.WithAllowUnknownColumns(true))
	if err != nil {
		log.Fatal(err)
	}

	dbscan.ErrNotFound = sql.ErrNoRows
	return &sqlxExecutor{
		scanner: scanner,
		db:      dbConn,
	}
}

func (e *sqlxExecutor) Exec(ctx context.Context, query IToSQL) (sql.Result, error) {
	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	return e.db.Exec(queryString, args...)
}

func (e *sqlxExecutor) Run(ctx context.Context, dest interface{}, query IToSQL) error {
	q, a, err := query.ToSql()
	if err != nil {
		return err
	}
	return e.selectByQuery(ctx, dest, q, a...)
}

func (e *sqlxExecutor) selectByQuery(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	var (
		val = reflect.ValueOf(dest)
	)
	if val.Kind() != reflect.Ptr {
		return errors.New("dest must be a pointer to a struct")
	}

	rows, err := e.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return err
	}

	if val.Elem().Kind() != reflect.Slice {
		err = e.scanner.ScanOne(dest, rows)
		return err
	}

	err = e.scanner.ScanAll(dest, rows)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}
