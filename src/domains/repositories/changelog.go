package repositories

import (
	"backend/src/domains/entities"
	"backend/src/modules/sql_executor"
	"context"

	"github.com/elgris/sqrl"
)

type changelogRepository struct {
	ICommonRepository
	executor sql_executor.ISQLExecutor
}

func NewChangelogRepository(executor sql_executor.ISQLExecutor) IChangelogRepository {
	return &changelogRepository{
		ICommonRepository: NewCommonRepository(),
		executor:          executor,
	}
}

func (r *changelogRepository) AddChangelogItems(ctx context.Context, items []*entities.ChangelogItem) error {
	q := sqrl.Insert(changelogTable).
		Columns(
			"target",
			"user_id",
			"table_id",
			"column_id",
			"row_id",
			"change",
			"changed_at",
		)

	for _, item := range items {
		q = q.Values(
			item.Target,
			item.UserID,
			item.TableID,
			item.ColumnID,
			item.RowID,
			item.Change,
			item.ChangedAt,
		)
	}

	q = q.PlaceholderFormat(sqrl.Dollar)

	_, err := r.executor.Exec(ctx, q)
	return err
}

func (r *changelogRepository) ListChangelogForCell(
	ctx context.Context,
	tableID string,
	columnID string,
	rowID int64,
) ([]*entities.ChangelogItemWithUserInfo, error) {
	q := sqrl.Select("*").
		From(changelogTableWithShortName).
		Join(usersTableWithShortName + " on cl.user_id = u.id").
		Where(sqrl.And{
			sqrl.Eq{"cl.target": entities.ChangeTargetCell},
			sqrl.Eq{"cl.table_id": tableID},
			sqrl.Eq{"cl.column_id": columnID},
			sqrl.Eq{"cl.row_id": rowID},
		}).
		PlaceholderFormat(sqrl.Dollar).
		OrderBy("changed_at ASC")

	var items []*entities.ChangelogItemWithUserInfo
	err := r.executor.Run(ctx, &items, q)
	return items, err
}

func (r *changelogRepository) ListChangelogForTable(
	ctx context.Context,
	tableID string,
) ([]*entities.ChangelogItemWithUserInfo, error) {
	q := sqrl.Select("*").
		From(changelogTableWithShortName).
		Join(usersTableWithShortName + " on cl.user_id = u.id").
		Where(sqrl.And{
			sqrl.Eq{"cl.target": entities.ChangeTargetTable},
			sqrl.Eq{"cl.table_id": tableID},
		}).
		PlaceholderFormat(sqrl.Dollar).
		OrderBy("changed_at ASC")

	var items []*entities.ChangelogItemWithUserInfo
	err := r.executor.Run(ctx, &items, q)
	return items, err
}
