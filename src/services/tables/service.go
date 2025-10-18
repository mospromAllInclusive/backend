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
	"github.com/xuri/excelize/v2"
)

const (
	tableIDTemplate  = "t_%s"
	columnIDTemplate = "col_%s"
)

type service struct {
	executor         sql_executor.ISQLExecutor
	repo             repositories.ITablesRepository
	changelogService services.IChangelogService
	fileService      services.IFileService
	keyMutex         key_mutex.IKeyMutex
}

func NewService(
	executor sql_executor.ISQLExecutor,
	repo repositories.ITablesRepository,
	changelogService services.IChangelogService,
	fileService services.IFileService,
) services.ITablesService {
	return &service{
		executor:         executor,
		repo:             repo,
		changelogService: changelogService,
		fileService:      fileService,
		keyMutex:         key_mutex.NewKeyMutex(),
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

	for _, col := range table.Columns {
		_, err = s.executor.Exec(ctx, table.CreateColumnIndexExpression(col.ID))
		if err != nil {
			return nil, err
		}
	}

	return s.repo.AddTable(ctx, table)
}

func (s *service) ImportTable(ctx context.Context, name string, databaseID int64, columns []string, data [][]*string) (*entities.Table, error) {
	rowsLimitPerInsert := 60000 / (len(columns) + 1)
	if rowsLimitPerInsert < 1 {
		return nil, fmt.Errorf("too many columns")
	}

	table := genDefaultTable(databaseID, name, columns)

	_, err := s.executor.Exec(ctx, table.CreateExpression())
	if err != nil {
		return nil, err
	}

	_, err = s.executor.Exec(ctx, table.CreateSortIndexExpression())
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(data); i += rowsLimitPerInsert {
		end := i + rowsLimitPerInsert
		if end > len(data) {
			end = len(data)
		}

		if err := s.repo.AddFullFilledRows(ctx, table, data[i:end]); err != nil {
			return nil, err
		}
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

	_, err = s.executor.Exec(ctx, table.CreateColumnIndexExpression(column.ID))
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateTable(ctx, table); err != nil {
		return nil, err
	}

	return table, nil
}

func (s *service) EditTableColumn(ctx context.Context, column *entities.TableColumn, tableID string) (*entities.Table, bool, error) {
	unlock := s.keyMutex.Lock(tableID)
	defer unlock()

	table, err := s.repo.GetTableByID(ctx, tableID, false)
	if err != nil {
		return nil, false, err
	}

	updated := false
	for _, col := range table.Columns {
		if col.ID == column.ID {
			if !col.NeedToBeUpdated(column) {
				break
			}
			col.Name = column.Name
			col.Type = column.Type
			col.Enum = column.Enum
			updated = true
			break
		}
	}

	if !updated {
		return table, false, nil
	}

	if err := s.repo.UpdateTable(ctx, table); err != nil {
		return nil, false, err
	}

	return table, true, nil
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

func (s *service) AddRow(ctx context.Context, userID int64, table *entities.Table, data map[string]*string, sortIndex *int64) (entities.TableRow, error) {
	newRow, err := s.repo.AddRow(ctx, table, data, sortIndex)
	if err != nil {
		return entities.TableRow{}, err
	}

	now := time.Now()
	changelog := make([]*entities.ChangelogItem, 0, len(data))
	for col, value := range data {
		if value == nil {
			continue
		}
		rawInfo := &entities.RawCellChangeInfo{
			Before:    nil,
			ChangedAt: now,
		}
		changelog = append(changelog, rawInfo.ToChangelogItem(userID, table.ID, newRow.GetID(), col, value))
	}

	if len(changelog) == 0 {
		return newRow, nil
	}

	return newRow, s.changelogService.WriteChangelog(ctx, changelog...)
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

func (s *service) SetCellValue(ctx context.Context, userID int64, tableID string, rowID int64, columnID string, value *string) error {
	rawChangeInfo, err := s.repo.SetCellValue(ctx, tableID, rowID, columnID, value)
	if err != nil {
		if s.repo.IsErrNoRows(err) {
			return nil
		}
		return err
	}

	return s.changelogService.WriteChangelog(ctx, rawChangeInfo.ToChangelogItem(userID, tableID, rowID, columnID, value))
}

func (s *service) ReadTable(ctx context.Context, table *entities.Table, params entities.ReadTableParams) ([]entities.TableRow, error) {
	return s.repo.ReadTable(ctx, table, &params)
}

func (s *service) GetTotalRows(ctx context.Context, table *entities.Table, params entities.ReadTableParams) (int64, error) {
	return s.repo.GetTotalRows(ctx, table, &params)
}

func (s *service) ExportTable(ctx context.Context, table *entities.Table) (*excelize.File, error) {
	rows, err := s.repo.ReadTable(ctx, table, nil)
	if err != nil {
		return nil, err
	}

	return s.fileService.CreateExcel(table, rows)
}

func genUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func genDefaultTable(databaseID int64, name string, columns []string) *entities.Table {
	table := &entities.Table{
		ID:         fmt.Sprintf(tableIDTemplate, genUUID()),
		Name:       name,
		DatabaseID: databaseID,
		Columns:    make([]*entities.TableColumn, 0, len(columns)),
	}

	for _, col := range columns {
		table.Columns = append(table.Columns, &entities.TableColumn{
			ID:   fmt.Sprintf(columnIDTemplate, genUUID()),
			Name: col,
			Type: entities.ColumnTypeText,
		})
	}

	return table
}
