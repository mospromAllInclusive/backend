package databases

import (
	"backend/src/domains/entities"
	"backend/src/handlers/common"
	"sort"
	"time"
)

type databaseResponse struct {
	ID        int64                 `json:"id"`
	Name      string                `json:"name"`
	Role      entities.Role         `json:"role,omitempty"`
	CreatedAt time.Time             `json:"created_at"`
	Tables    common.TablesResponse `json:"tables"`
}

func newDatabaseResponse(database *entities.Database) *databaseResponse {
	return &databaseResponse{
		ID:        database.ID,
		Name:      database.Name,
		CreatedAt: database.CreatedAt,
	}
}

type databaseListResponse []*databaseResponse

func newDatabaseListResponse(databases []*entities.UsersDatabase, tables []*entities.Table) databaseListResponse {
	tablesByDb := make(map[int64][]*entities.Table, len(databases))
	for _, table := range tables {
		tablesByDb[table.DatabaseID] = append(tablesByDb[table.DatabaseID], table)
	}
	res := make(databaseListResponse, 0, len(databases))
	for _, database := range databases {
		res = append(res, &databaseResponse{
			ID:        database.DatabaseID,
			Name:      database.Name,
			Role:      database.Role,
			CreatedAt: database.CreatedAt,
			Tables:    common.NewTablesResponse(tablesByDb[database.DatabaseID]),
		})
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})

	return res
}

type databasesUserResponse struct {
	*common.UserInfoResponse
	Role entities.Role `json:"role"`
}

func newDatabasesUserResponse(user *entities.DatabasesUser) *databasesUserResponse {
	return &databasesUserResponse{
		UserInfoResponse: common.NewUserInfoResponse(user.User),
		Role:             user.Role,
	}
}

type databaseUsersListResponse []*databasesUserResponse

func newDatabaseUsersListResponse(users []*entities.DatabasesUser) databaseUsersListResponse {
	res := make(databaseUsersListResponse, 0, len(users))
	for _, user := range users {
		res = append(res, newDatabasesUserResponse(user))
	}
	return res
}
