package repositories

import (
	"database/sql"
	"errors"
)

type commonRepository struct {
}

func NewCommonRepository() ICommonRepository {
	return &commonRepository{}
}

func (r *commonRepository) IsErrNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
