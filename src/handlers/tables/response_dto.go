package tables

import (
	"backend/src/domains/entities"
	"backend/src/handlers/common"
)

type rowResponse struct {
	ID   int64                  `json:"id"`
	Data map[string]interface{} `json:"data"`
}

func newRowResponse(row entities.TableRow) *rowResponse {
	res := &rowResponse{
		ID:   row.GetID(),
		Data: make(map[string]interface{}),
	}

	for k, v := range row {
		if k == "id" {
			continue
		}
		res.Data[k] = v
	}
	return res
}

type tableWithDataResponse struct {
	Table *common.TableResponse `json:"table"`
	Rows  []*rowResponse        `json:"rows"`
}

func newTableWithDataResponse(table *entities.Table, rows []entities.TableRow) *tableWithDataResponse {
	res := &tableWithDataResponse{
		Table: common.NewTableResponse(table),
		Rows:  make([]*rowResponse, 0, len(rows)),
	}

	for _, row := range rows {
		res.Rows = append(res.Rows, newRowResponse(row))
	}

	return res
}
