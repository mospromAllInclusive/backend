package repositories

import (
	"backend/src/domains/entities"
	"backend/src/modules/sql_executor"
	"context"
	"time"

	"github.com/elgris/sqrl"
)

type databasesRepository struct {
	ICommonRepository
	executor sql_executor.ISQLExecutor
}

func NewDatabasesRepository(executor sql_executor.ISQLExecutor) IDatabasesRepository {
	return &databasesRepository{
		ICommonRepository: NewCommonRepository(),
		executor:          executor,
	}
}

func (r *databasesRepository) AddDatabase(ctx context.Context, name string) (*entities.Database, error) {
	q := sqrl.Insert(databasesTable).
		Columns("name").
		Values(name).
		PlaceholderFormat(sqrl.Dollar).
		Returning("*")

	createdDBDatabase := &entities.Database{}
	err := r.executor.Run(ctx, createdDBDatabase, q)
	if err != nil {
		return nil, err
	}
	return createdDBDatabase, err
}

func (r *databasesRepository) UpsertUsersDatabase(ctx context.Context, usersDatabase *entities.UsersDatabase) (*entities.UsersDatabase, error) {
	q := sqrl.Insert(usersDatabasesTable).
		Columns("user_id, database_id, role").
		Values(usersDatabase.UserID, usersDatabase.DatabaseID, usersDatabase.Role).
		PlaceholderFormat(sqrl.Dollar).
		Suffix(`ON CONFLICT on constraint users_databases_pkey do update SET 
role = EXCLUDED.role,
deleted_at = null RETURNING *`)

	upsertedUsersDatabase := &entities.UsersDatabase{}
	err := r.executor.Run(ctx, upsertedUsersDatabase, q)
	if err != nil {
		return nil, err
	}
	return upsertedUsersDatabase, err
}

func (r *databasesRepository) DeleteUsersDatabaseRelation(ctx context.Context, userID, databaseID int64) error {
	q := sqrl.Update(usersDatabasesTable).
		Set("deleted_at", time.Now()).
		Where(sqrl.Eq{"user_id": userID, "database_id": databaseID}).
		PlaceholderFormat(sqrl.Dollar)

	_, err := r.executor.Exec(ctx, q)
	return err
}

func (r *databasesRepository) GetDatabaseByID(ctx context.Context, id int64) (*entities.Database, error) {
	q := sqrl.Select("*").
		From(databasesTable).
		Where(sqrl.Eq{"id": id, "deleted_at": nil}).
		PlaceholderFormat(sqrl.Dollar)

	dbDatabase := &entities.Database{}
	err := r.executor.Run(ctx, dbDatabase, q)
	if err != nil {
		return nil, err
	}
	return dbDatabase, err
}

func (r *databasesRepository) GetUsersDatabases(ctx context.Context, userID int64) ([]*entities.UsersDatabase, error) {
	q := sqrl.Select("udb.*, db.name").
		From(usersDatabasesTableWithShortName).
		Join(databasesTableWithShortName + " on udb.database_id = db.id").
		Where(sqrl.And{
			sqrl.Eq{"udb.user_id": userID},
			sqrl.Eq{"udb.deleted_at": nil},
			sqrl.Eq{"db.deleted_at": nil},
		}).
		PlaceholderFormat(sqrl.Dollar)

	var usersDatabases []*entities.UsersDatabase
	err := r.executor.Run(ctx, &usersDatabases, q)
	return usersDatabases, err
}

func (r *databasesRepository) GetDatabasesUsers(ctx context.Context, databaseID int64) ([]*entities.DatabasesUser, error) {
	q := sqrl.Select("u.*, udb.role").
		From(usersDatabasesTableWithShortName).
		Join(usersTableWithShortName + " on udb.user_id = u.id").
		Where(sqrl.And{
			sqrl.Eq{"udb.database_id": databaseID},
			sqrl.Eq{"udb.deleted_at": nil},
			sqrl.Eq{"u.deleted_at": nil},
		}).
		PlaceholderFormat(sqrl.Dollar)

	var databasesUsers []*entities.DatabasesUser
	err := r.executor.Run(ctx, &databasesUsers, q)
	return databasesUsers, err
}

func (r *databasesRepository) GetDatabasesUsersIDs(ctx context.Context, databaseID int64) ([]int64, error) {
	q := sqrl.Select("user_id").
		From(usersDatabasesTable).
		Where(sqrl.And{
			sqrl.Eq{"database_id": databaseID},
			sqrl.Eq{"deleted_at": nil},
		}).
		PlaceholderFormat(sqrl.Dollar)

	var ids []int64
	err := r.executor.Run(ctx, &ids, q)
	return ids, err
}

func (r *databasesRepository) GetUsersDatabaseRole(ctx context.Context, userID, databaseID int64) (entities.Role, error) {
	q := sqrl.Select("role").
		From(usersDatabasesTable).
		Where(sqrl.And{
			sqrl.Eq{"user_id": userID},
			sqrl.Eq{"database_id": databaseID},
			sqrl.Eq{"deleted_at": nil},
		}).
		PlaceholderFormat(sqrl.Dollar)

	var role entities.Role
	err := r.executor.Run(ctx, &role, q)
	return role, err
}
