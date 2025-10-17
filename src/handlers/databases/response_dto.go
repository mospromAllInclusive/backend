package databases

import (
	"backend/src/domains/entities"
	"sort"
	"time"
)

type databaseResponse struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	Role      entities.Role `json:"role,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
}

func newDatabaseResponse(database *entities.Database) *databaseResponse {
	return &databaseResponse{
		ID:        database.ID,
		Name:      database.Name,
		CreatedAt: database.CreatedAt,
	}
}

type databaseListResponse []*databaseResponse

func newDatabaseListResponse(databases []*entities.UsersDatabase) databaseListResponse {
	res := make(databaseListResponse, 0, len(databases))
	for _, database := range databases {
		res = append(res, &databaseResponse{
			ID:        database.DatabaseID,
			Name:      database.Name,
			Role:      database.Role,
			CreatedAt: database.CreatedAt,
		})
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})

	return res
}
