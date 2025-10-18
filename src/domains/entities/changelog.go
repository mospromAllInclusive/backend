package entities

import (
	"time"

	"github.com/AlekSi/pointer"
)

type ChangeTarget string

const (
	ChangeTargetCell     ChangeTarget = "cell"
	ChangeTargetTable    ChangeTarget = "table"
	ChangeTargetDatabase ChangeTarget = "database"
)

type ChangelogItem struct {
	ChangeID  int64         `db:"change_id"`
	Target    ChangeTarget  `db:"target"`
	UserID    int64         `db:"user_id"`
	TableID   *string       `db:"table_id"`
	ColumnID  *string       `db:"column_id"`
	RowID     *int64        `db:"row_id"`
	Change    JSONB[Change] `db:"change"`
	ChangedAt time.Time     `db:"changed_at"`
}

type ChangelogItemWithUserInfo struct {
	*ChangelogItem
	*User
}

type Change struct {
	CellChange *CellChange `json:"cell_change"`
}

type CellChange struct {
	Before *string `json:"before"`
	After  *string `json:"after"`
}

type RawCellChangeInfo struct {
	Before    *string   `db:"before"`
	ChangedAt time.Time `db:"changed_at"`
}

func (i *RawCellChangeInfo) ToChangelogItem(
	userID int64,
	tableID string,
	rowID int64,
	columnID string,
	value *string,
) *ChangelogItem {
	return &ChangelogItem{
		Target:   ChangeTargetCell,
		UserID:   userID,
		TableID:  pointer.To(tableID),
		ColumnID: pointer.To(columnID),
		RowID:    pointer.To(rowID),
		Change: JSONB[Change]{
			v: &Change{
				CellChange: &CellChange{
					Before: i.Before,
					After:  value,
				},
			},
		},
		ChangedAt: i.ChangedAt,
	}
}
