package users

import (
	"backend/src/domains/entities"
	"backend/src/domains/repositories"
	"backend/src/services"
	"context"
)

type service struct {
	repo repositories.IUsersRepository
}

func NewService(repo repositories.IUsersRepository) services.IUsersService {
	return &service{
		repo: repo,
	}
}

func (s *service) AddUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	return s.repo.CreateUser(ctx, user)
}

func (s *service) FindUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	res, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil && s.repo.IsErrNoRows(err) {
		return nil, ErrorUserNotFound{}
	}
	return res, err
}

func (s *service) FindUserByID(ctx context.Context, id int64) (*entities.User, error) {
	res, err := s.repo.GetUserByID(ctx, id)
	if err != nil && s.repo.IsErrNoRows(err) {
		return nil, ErrorUserNotFound{}
	}
	return res, err
}
