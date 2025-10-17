package repositories

import (
	"backend/src/domains/entities"
	"backend/src/modules/sql_executor"
	"context"
	"fmt"
	"time"

	"github.com/elgris/sqrl"
)

type tablesRepository struct {
	ICommonRepository
	executor sql_executor.ISQLExecutor
}

func NewTablesRepository(executor sql_executor.ISQLExecutor) ITablesRepository {
	return &tablesRepository{
		ICommonRepository: NewCommonRepository(),
		executor:          executor,
	}
}

func (r *tablesRepository) AddTable(ctx context.Context, table *entities.Table) (*entities.Table, error) {
	dbTable := table.ToDBTable()
	q := sqrl.Insert(tablesTable).
		Columns("id, name, database_id, columns").
		Values(dbTable.ID, dbTable.Name, dbTable.DatabaseID, dbTable.Columns).
		PlaceholderFormat(sqrl.Dollar).
		Returning("*")

	createdDBTable := &entities.DBTable{}
	err := r.executor.Run(ctx, createdDBTable, q)
	if err != nil {
		return nil, err
	}
	return createdDBTable.ToTable(), err
}

func (r *tablesRepository) DeleteTable(ctx context.Context, id string) error {
	q := sqrl.Update(tablesTable).
		Set("deleted_at", time.Now()).
		Where(sqrl.Eq{"id": id}).
		PlaceholderFormat(sqrl.Dollar)

	_, err := r.executor.Exec(ctx, q)
	return err
}

func (r *tablesRepository) RestoreTable(ctx context.Context, id string) error {
	q := sqrl.Update(tablesTable).
		Set("deleted_at", nil).
		Where(sqrl.Eq{"id": id}).
		PlaceholderFormat(sqrl.Dollar)

	_, err := r.executor.Exec(ctx, q)
	return err
}

func (r *tablesRepository) GetTableByID(ctx context.Context, id string, withDeleted bool) (*entities.Table, error) {
	q := sqrl.Select("*").
		From(tablesTable).
		Where(sqrl.Eq{"id": id}).
		PlaceholderFormat(sqrl.Dollar)

	if !withDeleted {
		q = q.Where(sqrl.Eq{"deleted_at": nil})
	}

	dbTable := &entities.DBTable{}
	err := r.executor.Run(ctx, dbTable, q)
	if err != nil {
		return nil, err
	}
	return dbTable.ToTable(), err
}

func (r *tablesRepository) UpdateTable(ctx context.Context, table *entities.Table) error {
	dbTable := table.ToDBTable()
	q := sqrl.Update(tablesTable).
		Set("name", dbTable.Name).
		Set("columns", dbTable.Columns).
		Set("database_id", dbTable.DatabaseID).
		Where(sqrl.Eq{"id": table.ID}).
		PlaceholderFormat(sqrl.Dollar)

	_, err := r.executor.Exec(ctx, q)
	return err
}

func (r *tablesRepository) ListByDatabaseID(ctx context.Context, databaseID int64) ([]*entities.Table, error) {
	q := sqrl.Select("*").
		From(tablesTable).
		Where(sqrl.Eq{"database_id": databaseID, "deleted_at": nil}).
		PlaceholderFormat(sqrl.Dollar)

	var dbTables []*entities.DBTable
	err := r.executor.Run(ctx, &dbTables, q)
	if err != nil {
		return nil, err
	}
	tables := make([]*entities.Table, len(dbTables))
	for i, dbTable := range dbTables {
		tables[i] = dbTable.ToTable()
	}
	return tables, err
}

func (r *tablesRepository) ListByDatabaseIDs(ctx context.Context, databaseIDs []int64) ([]*entities.Table, error) {
	q := sqrl.Select("*").
		From(tablesTable).
		Where(sqrl.Eq{"database_id": databaseIDs, "deleted_at": nil}).
		PlaceholderFormat(sqrl.Dollar)

	var dbTables []*entities.DBTable
	err := r.executor.Run(ctx, &dbTables, q)
	if err != nil {
		return nil, err
	}
	tables := make([]*entities.Table, len(dbTables))
	for i, dbTable := range dbTables {
		tables[i] = dbTable.ToTable()
	}
	return tables, err
}

func (r *tablesRepository) AddRow(ctx context.Context, table *entities.Table, data map[string]*string, sortIndex *int64) (entities.TableRow, error) {
	cols := []string{"sort_index_version"}
	values := []interface{}{time.Now().UnixNano()}
	if sortIndex != nil {
		cols = append(cols, "sort_index")
		values = append(values, *sortIndex)
	}

	for colID, value := range data {
		cols = append(cols, colID)
		values = append(values, value)
	}

	q := sqrl.Insert(fmt.Sprintf("%s.%s", entities.UsersTablespace, table.ID)).
		Columns(cols...).
		Values(values...).
		PlaceholderFormat(sqrl.Dollar).
		Returning(table.ReturningCols()...)

	row := make(entities.TableRow)
	err := r.executor.Run(ctx, &row, q)
	return row, err
}

func (r *tablesRepository) DeleteRow(ctx context.Context, tableID string, rowID int64) error {
	q := sqrl.Update(fmt.Sprintf("%s.%s", entities.UsersTablespace, tableID)).
		Set("deleted_at", time.Now()).
		Where(sqrl.Eq{"id": rowID}).
		PlaceholderFormat(sqrl.Dollar)

	_, err := r.executor.Exec(ctx, q)
	return err
}

func (r *tablesRepository) RestoreRow(ctx context.Context, tableID string, rowID int64) error {
	q := sqrl.Update(fmt.Sprintf("%s.%s", entities.UsersTablespace, tableID)).
		Set("deleted_at", nil).
		Where(sqrl.Eq{"id": rowID}).
		PlaceholderFormat(sqrl.Dollar)

	_, err := r.executor.Exec(ctx, q)
	return err
}

func (r *tablesRepository) MoveRow(ctx context.Context, tableID string, rowID int64, sortIndex int64) error {
	q := sqrl.Update(fmt.Sprintf("%s.%s", entities.UsersTablespace, tableID)).
		Set("sort_index", sortIndex).
		Set("sort_index_version", time.Now().UnixNano()).
		Where(sqrl.Eq{"id": rowID}).
		PlaceholderFormat(sqrl.Dollar)

	_, err := r.executor.Exec(ctx, q)
	return err
}

func (r *tablesRepository) SetCellValue(ctx context.Context, tableID string, rowID int64, columnID string, value *string) error {
	q := sqrl.Update(fmt.Sprintf("%s.%s", entities.UsersTablespace, tableID)).
		Set(columnID, value).
		Where(sqrl.Eq{"id": rowID}).
		PlaceholderFormat(sqrl.Dollar)

	_, err := r.executor.Exec(ctx, q)
	return err
}

func (r *tablesRepository) ReadTable(ctx context.Context, table *entities.Table) ([]entities.TableRow, error) {
	q := sqrl.Select(table.ReturningCols()...).
		From(fmt.Sprintf("%s.%s", entities.UsersTablespace, table.ID)).
		Where(sqrl.Eq{"deleted_at": nil}).
		OrderBy("sort_index ASC", "sort_index_version DESC").
		PlaceholderFormat(sqrl.Dollar)

	var rows []entities.TableRow
	err := r.executor.Run(ctx, &rows, q)
	return rows, err
}
