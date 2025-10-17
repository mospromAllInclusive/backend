package tables

import (
	"backend/src/domains/entities"
	"backend/src/domains/repositories"
	"backend/src/modules/key_mutex"
	"backend/src/modules/sql_executor"
	"backend/src/services"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/google/uuid"
)

const (
	tableIDTemplate  = "t_%s"
	columnIDTemplate = "col_%s"
)

type service struct {
	executor sql_executor.ISQLExecutor
	repo     repositories.ITablesRepository
	keyMutex key_mutex.IKeyMutex
}

func NewService(executor sql_executor.ISQLExecutor, repo repositories.ITablesRepository) services.ITablesService {
	return &service{
		executor: executor,
		repo:     repo,
		keyMutex: key_mutex.NewKeyMutex(),
	}
}

func (s *service) CreateTable(ctx context.Context, table *entities.Table) (*entities.Table, error) {
	table.ID = fmt.Sprintf(tableIDTemplate, genUUID())

	for _, col := range table.Columns {
		col.ID = fmt.Sprintf(columnIDTemplate, genUUID())
	}

	_, err := s.executor.Exec(ctx, table.CreateExpression())
	if err != nil {
		return nil, err
	}

	_, err = s.executor.Exec(ctx, table.CreateSortIndexExpression())
	if err != nil {
		return nil, err
	}

	return s.repo.AddTable(ctx, table)
}

func (s *service) DeleteTable(ctx context.Context, tableID string) error {
	return s.repo.DeleteTable(ctx, tableID)
}

func (s *service) RestoreTable(ctx context.Context, tableID string) error {
	return s.repo.RestoreTable(ctx, tableID)
}

func (s *service) AddColumnToTable(ctx context.Context, column *entities.TableColumn, tableID string) (*entities.Table, error) {
	unlock := s.keyMutex.Lock(tableID)
	defer unlock()

	table, err := s.repo.GetTableByID(ctx, tableID, false)
	if err != nil {
		return nil, err
	}

	column.ID = fmt.Sprintf(columnIDTemplate, genUUID())
	table.Columns = append(table.Columns, column)

	_, err = s.executor.Exec(ctx, table.AddColumnExpression(column))
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateTable(ctx, table); err != nil {
		return nil, err
	}

	return table, nil
}

func (s *service) DeleteColumn(ctx context.Context, columnID string, tableID string) (*entities.Table, error) {
	unlock := s.keyMutex.Lock(tableID)
	defer unlock()

	table, err := s.repo.GetTableByID(ctx, tableID, false)
	if err != nil {
		return nil, err
	}

	found := false
	for _, col := range table.Columns {
		if col.ID == columnID {
			col.DeletedAt = pointer.To(time.Now())
			found = true
		}
	}

	if !found {
		return nil, ErrorColumnNotFound{}
	}

	if err := s.repo.UpdateTable(ctx, table); err != nil {
		return nil, err
	}

	return table, nil
}

func (s *service) RestoreColumn(ctx context.Context, columnID string, tableID string) (*entities.Table, error) {
	unlock := s.keyMutex.Lock(tableID)
	defer unlock()

	table, err := s.repo.GetTableByID(ctx, tableID, false)
	if err != nil {
		return nil, err
	}

	found := false
	for _, col := range table.Columns {
		if col.ID == columnID {
			col.DeletedAt = nil
			found = true
		}
	}

	if !found {
		return nil, ErrorColumnNotFound{}
	}

	if err := s.repo.UpdateTable(ctx, table); err != nil {
		return nil, err
	}

	return table, nil
}

func (s *service) GetTableByID(ctx context.Context, id string, withDeleted bool) (*entities.Table, error) {
	table, err := s.repo.GetTableByID(ctx, id, withDeleted)
	if err != nil && s.repo.IsErrNoRows(err) {
		return nil, ErrorTableNotFound{}
	}
	return table, err
}

func (s *service) ListByDatabaseID(ctx context.Context, databaseID int64) ([]*entities.Table, error) {
	return s.repo.ListByDatabaseID(ctx, databaseID)
}

func (s *service) ListByDatabaseIDs(ctx context.Context, databaseIDs []int64) ([]*entities.Table, error) {
	return s.repo.ListByDatabaseIDs(ctx, databaseIDs)
}

func (s *service) AddRow(ctx context.Context, table *entities.Table, sortIndex *int64) (entities.TableRow, error) {
	return s.repo.AddRow(ctx, table, sortIndex)
}

func (s *service) DeleteRow(ctx context.Context, tableID string, rowID int64) error {
	return s.repo.DeleteRow(ctx, tableID, rowID)
}

func (s *service) RestoreRow(ctx context.Context, tableID string, rowID int64) error {
	return s.repo.RestoreRow(ctx, tableID, rowID)
}

func (s *service) MoveRow(ctx context.Context, tableID string, rowID int64, sortIndex int64) error {
	return s.repo.MoveRow(ctx, tableID, rowID, sortIndex)
}

func (s *service) SetCellValue(ctx context.Context, tableID string, rowID int64, columnID string, value *string) error {
	return s.repo.SetCellValue(ctx, tableID, rowID, columnID, value)
}

func (s *service) ReadTable(ctx context.Context, table *entities.Table) ([]entities.TableRow, error) {
	return s.repo.ReadTable(ctx, table)
}

func genUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
