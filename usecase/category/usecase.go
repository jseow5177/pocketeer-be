package category

import (
	"context"

	"github.com/jseow5177/pockteer-be/entity"
)

type CategoryUseCase struct {
	Repo
}

func NewCategoryUseCase(repo Repo) *CategoryUseCase {
	return &CategoryUseCase{
		repo,
	}
}

func (uc *CategoryUseCase) CreateCategory(ctx context.Context, c *entity.Category) error {
	return uc.Create(ctx, c)
}

func (uc *CategoryUseCase) GetCategory(ctx context.Context, f *entity.CategoryFilter) (*entity.Category, error) {
	return nil, nil
}
