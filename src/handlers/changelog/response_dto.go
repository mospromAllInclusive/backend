package changelog

import (
	"backend/src/domains/entities"
	"backend/src/handlers/common"
	"time"

	"github.com/AlekSi/pointer"
)

type cellChangelogItemResponse struct {
	ChangeID  int64                    `json:"change_id"`
	Before    *string                  `json:"before"`
	After     *string                  `json:"after"`
	ChangedAt time.Time                `json:"changed_at"`
	User      *common.UserInfoResponse `json:"user"`
}

func newCellChangelogItemResponse(item *entities.ChangelogItemWithUserInfo) *cellChangelogItemResponse {
	return &cellChangelogItemResponse{
		ChangeID:  item.ChangeID,
		Before:    item.Change.Get().CellChange.Before,
		After:     item.Change.Get().CellChange.After,
		ChangedAt: item.ChangedAt,
		User:      common.NewUserInfoResponse(item.User),
	}
}

type cellChangelogResponse []*cellChangelogItemResponse

func newCellChangelogResponse(changelog []*entities.ChangelogItemWithUserInfo) cellChangelogResponse {
	res := make(cellChangelogResponse, 0, len(changelog))
	for _, item := range changelog {
		res = append(res, newCellChangelogItemResponse(item))
	}
	return res
}

type tableChangelogResponse []*tableChangelogItemResponse

func newTableChangelogResponse(changelog []*entities.ChangelogItemWithUserInfo) tableChangelogResponse {
	res := make(tableChangelogResponse, 0, len(changelog))
	for _, item := range changelog {
		if logItem := newTableChangelogItemResponse(item); logItem != nil {
			res = append(res, logItem)
		}
	}
	return res
}

type tableChangelogItemResponse struct {
	ChangeID      int64                     `json:"change_id"`
	ChangedEntity entities.ChangedEntity    `json:"changed_entity"`
	ChangeType    entities.ChangeType       `json:"change_type"`
	BeforeColumn  *common.ColumnForResponse `json:"before_column"`
	AfterColumn   *common.ColumnForResponse `json:"after_column"`
	BeforeRow     rowForChangelog           `json:"before_row"`
	AfterRow      rowForChangelog           `json:"after_row"`
	ChangedAt     time.Time                 `json:"changed_at"`
	User          *common.UserInfoResponse  `json:"user"`
}

func newTableChangelogItemResponse(item *entities.ChangelogItemWithUserInfo) *tableChangelogItemResponse {
	change := item.Change.Get()
	res := &tableChangelogItemResponse{
		ChangeID:      item.ChangeID,
		ChangedEntity: change.ChangedEntity,
		ChangedAt:     item.ChangedAt,
		User:          common.NewUserInfoResponse(item.User),
	}

	switch change.ChangedEntity {
	case entities.ChangedEntityRow:
		res.ChangeType = change.RowChange.ChangeType
		res.BeforeRow = newRowForChangelog(change.RowChange.Before)
		res.AfterRow = newRowForChangelog(change.RowChange.After)
	case entities.ChangedEntityColumn:
		res.ChangeType = change.ColumnChange.ChangeType
		if change.ColumnChange.Before != nil {
			res.BeforeColumn = pointer.To(common.NewColumnForResponse(change.ColumnChange.Before))
		}
		if change.ColumnChange.After != nil {
			res.AfterColumn = pointer.To(common.NewColumnForResponse(change.ColumnChange.After))
		}

	default:
		return nil
	}

	return res
}

type rowForChangelog []*rowItemInfoForChangelog

type rowItemInfoForChangelog struct {
	ColumnID   string
	ColumnName string
	Value      any
}

func newRowForChangelog(row entities.RowInfoForChangelog) rowForChangelog {
	if len(row) == 0 {
		return nil
	}

	res := make(rowForChangelog, 0, len(row))
	for _, item := range row {
		res = append(res, &rowItemInfoForChangelog{
			ColumnID:   item.ColumnID,
			ColumnName: item.ColumnName,
			Value:      item.Value,
		})
	}
	return res
}
