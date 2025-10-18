package services

import (
	"backend/src/domains/entities"
	"context"
	"encoding/csv"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type IUsersService interface {
	AddUser(ctx context.Context, user *entities.User) (*entities.User, error)
	FindUserByEmail(ctx context.Context, email string) (*entities.User, error)
	FindUserByID(ctx context.Context, id int64) (*entities.User, error)
	ListUsers(ctx context.Context) ([]*entities.User, error)
}

type IAuthService interface {
	JWTAuthMiddleware() gin.HandlerFunc
	Login(ctx context.Context, email string, password string) (*entities.User, string, error)
}

type ITablesService interface {
	CreateTable(ctx context.Context, table *entities.Table) (*entities.Table, error)
	ImportTable(ctx context.Context, name string, databaseID int64, columns []string, data [][]*string) (*entities.Table, error)
	DeleteTable(ctx context.Context, id string) error
	RestoreTable(ctx context.Context, id string) error
	AddColumnToTable(ctx context.Context, column *entities.TableColumn, tableID string) (*entities.Table, error)
	EditTableColumn(ctx context.Context, column *entities.TableColumn, tableID string) (*entities.Table, error)
	DeleteColumn(ctx context.Context, columnID string, tableID string) (*entities.Table, error)
	RestoreColumn(ctx context.Context, columnID string, tableID string) (*entities.Table, error)
	GetTableByID(ctx context.Context, id string, withDeleted bool) (*entities.Table, error)
	ListByDatabaseID(ctx context.Context, databaseID int64) ([]*entities.Table, error)
	ListByDatabaseIDs(ctx context.Context, databaseIDs []int64) ([]*entities.Table, error)
	AddRow(ctx context.Context, userID int64, table *entities.Table, data map[string]*string, sortIndex *int64) (entities.TableRow, error)
	DeleteRow(ctx context.Context, tableID string, rowID int64) error
	RestoreRow(ctx context.Context, tableID string, rowID int64) error
	MoveRow(ctx context.Context, tableID string, rowID int64, sortIndex int64) error
	SetCellValue(ctx context.Context, userID int64, tableID string, rowID int64, columnID string, value *string) error
	ReadTable(ctx context.Context, table *entities.Table, params entities.ReadTableParams) ([]entities.TableRow, error)
	ExportTable(ctx context.Context, table *entities.Table) (*excelize.File, error)
}

type IDatabasesService interface {
	AddDatabase(ctx context.Context, userID int64, name string) (*entities.Database, error)
	UpsertUsersDatabase(ctx context.Context, usersDatabase *entities.UsersDatabase) (*entities.UsersDatabase, error)
	DeleteUsersDatabaseRelation(ctx context.Context, userID, databaseID int64) error
	GetUsersDatabases(ctx context.Context, userID int64) ([]*entities.UsersDatabase, error)
	GetDatabasesUsers(ctx context.Context, databaseID int64) ([]*entities.DatabasesUser, error)
	CheckUserRole(ctx context.Context, userID, databaseID int64, requiredRole entities.Role) (bool, error)
	GetUsersDatabaseRole(ctx context.Context, userID, databaseID int64) (entities.Role, error)
}

type IChangelogService interface {
	WriteChangelog(ctx context.Context, items ...*entities.ChangelogItem) error
	ListChangelogForCell(
		ctx context.Context,
		tableID string,
		columnID string,
		rowID int64,
	) ([]*entities.ChangelogItemWithUserInfo, error)
}

type IFileService interface {
	ReadFile(file *multipart.FileHeader) ([]string, [][]*string, error)
	ReadExcel(f *excelize.File) ([]string, [][]*string, error)
	ReadCSV(r *csv.Reader) ([]string, [][]*string, error)
	CreateExcel(table *entities.Table, data []entities.TableRow) (f *excelize.File, err error)
}
