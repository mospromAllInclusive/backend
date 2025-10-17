package resources

import (
	"backend/src/modules/sql_executor"
	"os"
)

type Resources struct {
	PostgresExecutor sql_executor.ISQLExecutor
}

func NewResources() *Resources {
	r := &Resources{}

	r.PostgresExecutor = sql_executor.NewSQlXExecutor(
		"postgres",
		os.Getenv("POSTGRES_URL"),
	)

	return r
}
