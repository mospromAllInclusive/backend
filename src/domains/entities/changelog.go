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

type ChangedEntity string

const (
	ChangedEntityCell   ChangedEntity = "cell"
	ChangedEntityRow    ChangedEntity = "row"
	ChangedEntityColumn ChangedEntity = "column"
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
	ChangedEntity ChangedEntity `json:"changed_entity"`
	CellChange    *CellChange   `json:"cell_change"`
	ColumnChange  *ColumnChange `json:"column_change"`
	RowChange     *RowChange    `json:"row_change"`
}

type CellChange struct {
	Before *string `json:"before"`
	After  *string `json:"after"`
}

type ChangeType string

const (
	ChangeTypeAdd    ChangeType = "add"
	ChangeTypeUpdate ChangeType = "update"
	ChangeTypeDelete ChangeType = "delete"
)

type ColumnChange struct {
	ChangeType ChangeType   `json:"change_type"`
	Before     *TableColumn `json:"before"`
	After      *TableColumn `json:"after"`
}

func (i *ColumnChange) ToChangelogItem(
	userID int64,
	tableID string,
	columnID string,
) *ChangelogItem {
	return &ChangelogItem{
		Target:   ChangeTargetTable,
		UserID:   userID,
		TableID:  pointer.To(tableID),
		ColumnID: pointer.To(columnID),
		Change: JSONB[Change]{
			v: &Change{
				ChangedEntity: ChangedEntityColumn,
				ColumnChange:  i,
			},
		},
		ChangedAt: time.Now(),
	}
}

type RowChange struct {
	ChangeType ChangeType          `json:"change_type"`
	Before     RowInfoForChangelog `json:"before"`
	After      RowInfoForChangelog `json:"after"`
}

func (i *RowChange) ToChangelogItem(
	userID int64,
	tableID string,
	rowID int64,
) *ChangelogItem {
	return &ChangelogItem{
		Target:  ChangeTargetTable,
		UserID:  userID,
		TableID: pointer.To(tableID),
		RowID:   pointer.To(rowID),
		Change: JSONB[Change]{
			v: &Change{
				ChangedEntity: ChangedEntityRow,
				RowChange:     i,
			},
		},
		ChangedAt: time.Now(),
	}
}

type RowInfoForChangelog []*RowItemForChangelog

type RowItemForChangelog struct {
	ColumnID   string
	ColumnName string
	Value      any
}

func NewRowInfoForChangelog(table *Table, row TableRow) RowInfoForChangelog {
	res := make(RowInfoForChangelog, 0, len(table.Columns))
	for _, column := range table.Columns {
		if column.DeletedAt != nil {
			continue
		}

		res = append(res, &RowItemForChangelog{
			ColumnID:   column.ID,
			ColumnName: column.Name,
			Value:      row[column.ID],
		})
	}

	return res
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
				ChangedEntity: ChangedEntityCell,
				CellChange: &CellChange{
					Before: i.Before,
					After:  value,
				},
			},
		},
		ChangedAt: i.ChangedAt,
	}
}
