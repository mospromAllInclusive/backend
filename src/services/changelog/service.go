package changelog

import (
	"backend/src/domains/entities"
	"backend/src/domains/repositories"
	"backend/src/services"
	"context"
)

type service struct {
	repo repositories.IChangelogRepository
}

func NewService(repo repositories.IChangelogRepository) services.IChangelogService {
	return &service{
		repo: repo,
	}
}

func (s *service) WriteChangelog(ctx context.Context, items ...*entities.ChangelogItem) error {
	return s.repo.AddChangelogItems(ctx, items)
}

func (s *service) ListChangelogForCell(
	ctx context.Context,
	tableID string,
	columnID string,
	rowID int64,
) ([]*entities.ChangelogItemWithUserInfo, error) {
	return s.repo.ListChangelogForCell(ctx, tableID, columnID, rowID)
}

func (s *service) ListChangelogForTable(
	ctx context.Context,
	tableID string,
) ([]*entities.ChangelogItemWithUserInfo, error) {
	return s.repo.ListChangelogForTable(ctx, tableID)
}
