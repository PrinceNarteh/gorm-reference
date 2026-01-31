package repository

import (
	"context"

	"gorm.io/gorm"
)

var _ QueryRepository = (*queryRepository)(nil)

type QueryRepository interface {
	Association(ctx context.Context, column string) *gorm.Association
}

type queryRepository struct {
	db *gorm.DB
}

func (r *queryRepository) Association(ctx context.Context, column string) *gorm.Association {
	return r.db.WithContext(ctx).Association(column)
}
