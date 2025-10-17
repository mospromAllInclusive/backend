package repositories

import (
	"backend/src/domains/entities"
	"context"
)

const (
	usersTable                       = "app.users"
	usersTableWithShortName          = "app.users as u"
	tablesTable                      = "app.tables"
	tablesTableWithShortName         = "app.tables as t"
	databasesTable                   = "app.databases"
	databasesTableWithShortName      = "app.databases as db"
	usersDatabasesTable              = "app.users_databases"
	usersDatabasesTableWithShortName = "app.users_databases as udb"
)

type ICommonRepository interface {
	IsErrNoRows(err error) bool
}

type IUsersRepository interface {
	ICommonRepository
	CreateUser(ctx context.Context, user *entities.User) (*entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	GetUserByID(ctx context.Context, id int64) (*entities.User, error)
	ListUsers(ctx context.Context) ([]*entities.User, error)
}

type ITablesRepository interface {
	ICommonRepository
	AddTable(ctx context.Context, table *entities.Table) (*entities.Table, error)
	DeleteTable(ctx context.Context, id string) error
	RestoreTable(ctx context.Context, id string) error
	GetTableByID(ctx context.Context, id string, withDeleted bool) (*entities.Table, error)
	UpdateTable(ctx context.Context, table *entities.Table) error
	ListByDatabaseID(ctx context.Context, databaseID int64) ([]*entities.Table, error)
	ListByDatabaseIDs(ctx context.Context, databaseIDs []int64) ([]*entities.Table, error)
	AddRow(ctx context.Context, table *entities.Table, data map[string]*string, sortIndex *int64) (entities.TableRow, error)
	DeleteRow(ctx context.Context, tableID string, rowID int64) error
	RestoreRow(ctx context.Context, tableID string, rowID int64) error
	MoveRow(ctx context.Context, tableID string, rowID int64, sortIndex int64) error
	SetCellValue(ctx context.Context, tableID string, rowID int64, columnID string, value *string) error
	ReadTable(ctx context.Context, table *entities.Table) ([]entities.TableRow, error)
}

type IDatabasesRepository interface {
	ICommonRepository
	AddDatabase(ctx context.Context, name string) (*entities.Database, error)
	UpsertUsersDatabase(ctx context.Context, usersDatabase *entities.UsersDatabase) (*entities.UsersDatabase, error)
	DeleteUsersDatabaseRelation(ctx context.Context, userID, databaseID int64) error
	GetDatabaseByID(ctx context.Context, id int64) (*entities.Database, error)
	GetUsersDatabases(ctx context.Context, userID int64) ([]*entities.UsersDatabase, error)
	GetDatabasesUsers(ctx context.Context, databaseID int64) ([]*entities.DatabasesUser, error)
	GetUsersDatabaseRole(ctx context.Context, userID, databaseID int64) (entities.Role, error)
}
