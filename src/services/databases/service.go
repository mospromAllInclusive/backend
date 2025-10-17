package databases

import (
	"backend/src/domains/entities"
	"backend/src/domains/repositories"
	"backend/src/services"
	"context"
)

type service struct {
	repo repositories.IDatabasesRepository
}

func NewService(repo repositories.IDatabasesRepository) services.IDatabasesService {
	return &service{
		repo: repo,
	}
}

func (s *service) AddDatabase(ctx context.Context, userID int64, name string) (*entities.Database, error) {
	createdDatabase, err := s.repo.AddDatabase(ctx, name)
	if err != nil {
		return nil, err
	}

	usersDatabase := &entities.UsersDatabase{
		UserID:     userID,
		DatabaseID: createdDatabase.ID,
		Role:       entities.RoleAdmin,
	}

	_, err = s.UpsertUsersDatabase(ctx, usersDatabase)
	if err != nil {
		return nil, err
	}

	return createdDatabase, nil
}

func (s *service) UpsertUsersDatabase(ctx context.Context, usersDatabase *entities.UsersDatabase) (*entities.UsersDatabase, error) {
	return s.repo.UpsertUsersDatabase(ctx, usersDatabase)
}

func (s *service) DeleteUsersDatabaseRelation(ctx context.Context, userID, databaseID int64) error {
	return s.repo.DeleteUsersDatabaseRelation(ctx, userID, databaseID)
}

func (s *service) GetUsersDatabases(ctx context.Context, userID int64) ([]*entities.UsersDatabase, error) {
	return s.repo.GetUsersDatabases(ctx, userID)
}

func (s *service) GetDatabasesUsers(ctx context.Context, databaseID int64) ([]*entities.DatabasesUser, error) {
	return s.repo.GetDatabasesUsers(ctx, databaseID)
}

func (s *service) CheckUserRole(ctx context.Context, userID, databaseID int64, requiredRole entities.Role) (bool, error) {
	role, err := s.repo.GetUsersDatabaseRole(ctx, userID, databaseID)
	if err != nil {
		if s.repo.IsErrNoRows(err) {
			return false, nil
		}
		return false, err
	}

	return role.Authorize(requiredRole), nil
}
