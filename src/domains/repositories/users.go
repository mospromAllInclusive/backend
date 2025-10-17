package repositories

import (
	"backend/src/domains/entities"
	"backend/src/modules/sql_executor"
	"context"

	"github.com/elgris/sqrl"
)

type usersRepository struct {
	ICommonRepository
	executor sql_executor.ISQLExecutor
}

func NewUsersRepository(executor sql_executor.ISQLExecutor) IUsersRepository {
	return &usersRepository{
		ICommonRepository: NewCommonRepository(),
		executor:          executor,
	}
}

func (r *usersRepository) CreateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	q := sqrl.Insert(usersTable).
		Columns("name, email, password").
		Values(user.Name, user.Email, user.Password).
		PlaceholderFormat(sqrl.Dollar).
		Returning("*")

	var task entities.User
	err := r.executor.Run(ctx, &task, q)
	return &task, err
}

func (r *usersRepository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	q := sqrl.Select("*").
		From(usersTable).
		Where(sqrl.Eq{"email": email, "deleted_at": nil}).
		PlaceholderFormat(sqrl.Dollar)

	var user entities.User
	err := r.executor.Run(ctx, &user, q)
	return &user, err
}

func (r *usersRepository) GetUserByID(ctx context.Context, id int64) (*entities.User, error) {
	q := sqrl.Select("*").
		From(usersTable).
		Where(sqrl.Eq{"id": id, "deleted_at": nil}).
		PlaceholderFormat(sqrl.Dollar)

	var user entities.User
	err := r.executor.Run(ctx, &user, q)
	return &user, err
}
