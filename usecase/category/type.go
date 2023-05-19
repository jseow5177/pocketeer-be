package category

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
)

type UseCase interface {
	GetCategory(ctx context.Context, userID string, req *presenter.GetCategoryRequest) (*entity.Category, error)
	GetCategories(ctx context.Context, userID string, req *presenter.GetCategoriesRequest) ([]*entity.Category, error)

	CreateCategory(ctx context.Context, userID string, req *presenter.CreateCategoryRequest) (*entity.Category, error)
	UpdateCategory(ctx context.Context, userID string, req *presenter.UpdateCategoryRequest) (*entity.Category, error)
}
