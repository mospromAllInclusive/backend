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
	Columns    []ColumnForResponse `json:"columns"`
	CreatedAt  time.Time           `json:"created_at"`
	TotalRows  *int64              `json:"total_rows,omitempty"`
}

type ColumnForResponse struct {
	Name string              `json:"name"`
	Type entities.ColumnType `json:"type"`
	ID   string              `json:"id"`
	Enum []string            `json:"enum"`
}

func NewTableResponse(table *entities.Table) *TableResponse {
	cols := make([]ColumnForResponse, 0, len(table.Columns))
	for _, col := range table.Columns {
		if col.DeletedAt != nil {
			continue
		}
		cols = append(cols, NewColumnForResponse(col))
	}
	return &TableResponse{
		ID:         table.ID,
		Name:       table.Name,
		DatabaseID: table.DatabaseID,
		Columns:    cols,
		CreatedAt:  table.CreatedAt,
	}
}

func NewColumnForResponse(col *entities.TableColumn) ColumnForResponse {
	return ColumnForResponse{
		Name: col.Name,
		Type: col.Type,
		ID:   col.ID,
		Enum: col.Enum,
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
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUserInfoResponse(user *entities.User) *UserInfoResponse {
	return &UserInfoResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}
