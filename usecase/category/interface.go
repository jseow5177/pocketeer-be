package category

import (
	"context"

	"github.com/jseow5177/pockteer-be/data/entity"
)

type UseCase interface {
	GetCategory(ctx context.Context, f *entity.CategoryFilter) (*entity.Category, error)
	GetCategories(ctx context.Context, f *entity.CategoryFilter) ([]*entity.Category, error)

	CreateCategory(ctx context.Context, c *entity.Category) error
	UpdateCategory(ctx context.Context, c *entity.Category) error
}
