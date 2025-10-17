package common

import (
	"backend/src/domains/entities"
	"sort"
	"time"
)

type TableResponse struct {
	ID         string              `json:"id"`
	Name       string              `json:"name"`
	DatabaseID int64               `json:"database_id"`
	Columns    []columnForResponse `json:"columns"`
	CreatedAt  time.Time           `json:"created_at"`
}

type columnForResponse struct {
	Name string              `json:"name"`
	Type entities.ColumnType `json:"type"`
	ID   string              `json:"id"`
}

func NewTableResponse(table *entities.Table) *TableResponse {
	cols := make([]columnForResponse, 0, len(table.Columns))
	for _, col := range table.Columns {
		cols = append(cols, columnForResponse{
			Name: col.Name,
			Type: col.Type,
			ID:   col.ID,
		})
	}
	return &TableResponse{
		ID:         table.ID,
		Name:       table.Name,
		DatabaseID: table.DatabaseID,
		Columns:    cols,
		CreatedAt:  table.CreatedAt,
	}
}

type TablesResponse []*TableResponse

func NewTablesResponse(tables []*entities.Table) TablesResponse {
	res := make(TablesResponse, 0, len(tables))
	for _, table := range tables {
		res = append(res, NewTableResponse(table))
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})

	return res
}

type UserInfoResponse struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUserInfoResponse(user *entities.User) *UserInfoResponse {
	return &UserInfoResponse{
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}
