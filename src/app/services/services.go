package services

import (
	"backend/src/app/repositories"
	"backend/src/app/resources"
	"backend/src/services"
	"backend/src/services/auth"
	"backend/src/services/changelog"
	"backend/src/services/databases"
	"backend/src/services/file_service"
	"backend/src/services/tables"
	"backend/src/services/users"
)

type Services struct {
	UsersService     services.IUsersService
	ChangelogService services.IChangelogService
	TablesService    services.ITablesService
	AuthService      services.IAuthService
	DatabasesService services.IDatabasesService
	FileService      services.IFileService
}

func NewServices(repos *repositories.Repositories, res *resources.Resources) *Services {
	s := &Services{}

	s.FileService = file_service.NewService()
	s.UsersService = users.NewService(repos.UsersRepository)
	s.ChangelogService = changelog.NewService(repos.ChangelogRepository)
	s.TablesService = tables.NewService(res.PostgresExecutor, repos.TablesRepository, s.ChangelogService, s.FileService)
	s.AuthService = auth.NewService(s.UsersService)
	s.DatabasesService = databases.NewService(repos.DatabasesRepository)

	return s
}
