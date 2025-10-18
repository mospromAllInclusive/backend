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
	WSHub            *web_sockets.Hub
}

func NewResources() *Resources {
	r := &Resources{}

	r.Ctx, r.Cancel = context.WithCancel(context.Background())

	r.PostgresExecutor = sql_executor.NewSQlXExecutor(
		"postgres",
		os.Getenv("POSTGRES_URL"),
	)

	r.WSHub = web_sockets.NewHub(r.Ctx)

	return r
}
