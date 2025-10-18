package resources

import (
	"backend/src/modules/sql_executor"
	"backend/src/modules/web_sockets"
	"context"
	"os"
)

type Resources struct {
	Ctx              context.Context
	Cancel           context.CancelFunc
	PostgresExecutor sql_executor.ISQLExecutor
	TablesWSHub      *web_sockets.Hub
	UsersWSHub       *web_sockets.Hub
}

func NewResources() *Resources {
	r := &Resources{}

	r.Ctx, r.Cancel = context.WithCancel(context.Background())

	r.PostgresExecutor = sql_executor.NewSQlXExecutor(
		"postgres",
		os.Getenv("POSTGRES_URL"),
	)

	r.TablesWSHub = web_sockets.NewHub(r.Ctx)
	r.UsersWSHub = web_sockets.NewHub(r.Ctx)

	return r
}
