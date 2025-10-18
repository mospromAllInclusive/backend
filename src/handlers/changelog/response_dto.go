package changelog

import (
	"backend/src/domains/entities"
	"backend/src/handlers/common"
	"time"
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
