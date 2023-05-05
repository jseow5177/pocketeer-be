package category

import (
	"context"

	"github.com/jseow5177/pockteer-be/entity"
)

type Reader interface {
	Get(ctx context.Context, f *entity.CategoryFilter) (*entity.Category, error)
	GetMany(ctx context.Context, f *entity.CategoryFilter) ([]*entity.Category, error)
}

type Writer interface {
	Create(ctx context.Context, c *entity.Category) error
	Update(ctx context.Context, c *entity.Category) error
}

type Repo interface {
	Reader
	Writer
}

type UseCase interface {
	GetCategory(ctx context.Context, f *entity.CategoryFilter) (*entity.Category, error)
	GetCategories(ctx context.Context, f *entity.CategoryFilter) ([]*entity.Category, error)

	CreateCategory(ctx context.Context, c *entity.Category) error
	UpdateCategory(ctx context.Context, c *entity.Category) error
}
