package repositories

import (
	"backend/src/app/resources"
	"backend/src/domains/repositories"
)

type Repositories struct {
	UsersRepository     repositories.IUsersRepository
	TablesRepository    repositories.ITablesRepository
	DatabasesRepository repositories.IDatabasesRepository
	ChangelogRepository repositories.IChangelogRepository
}

func NewRepositories(res *resources.Resources) *Repositories {
	r := &Repositories{}

	r.UsersRepository = repositories.NewUsersRepository(res.PostgresExecutor)
	r.TablesRepository = repositories.NewTablesRepository(res.PostgresExecutor)
	r.DatabasesRepository = repositories.NewDatabasesRepository(res.PostgresExecutor)
	r.ChangelogRepository = repositories.NewChangelogRepository(res.PostgresExecutor)

	return r
}
