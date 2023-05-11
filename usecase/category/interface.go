package category

import (
	"context"

	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/data/presenter"
)

type UseCase interface {
	GetCategory(ctx context.Context, req *presenter.GetCategoryRequest) (*entity.Category, error)
	GetCategories(ctx context.Context, userID string, req *presenter.GetCategoriesRequest) ([]*entity.Category, error)

	CreateCategory(ctx context.Context, userID string, req *presenter.CreateCategoryRequest) (*entity.Category, error)
	UpdateCategory(ctx context.Context, req *presenter.UpdateCategoryRequest) (*entity.Category, error)
}
