package category

import (
	"context"
	"time"

	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

type CategoryUseCase struct {
	categoryRepo repo.CategoryRepo
}

func NewCategoryUseCase(categoryRepo repo.CategoryRepo) UseCase {
	return &CategoryUseCase{
		categoryRepo: categoryRepo,
	}
}

func (uc *CategoryUseCase) CreateCategory(ctx context.Context, c *entity.Category) error {
	now := time.Now().Unix()

	c.CreateTime = goutil.Uint64(uint64(now))
	c.UpdateTime = goutil.Uint64(uint64(now))

	return uc.categoryRepo.Create(ctx, c)
}

func (uc *CategoryUseCase) UpdateCategory(ctx context.Context, c *entity.Category) error {
	return nil
}

func (uc *CategoryUseCase) GetCategory(ctx context.Context, f *entity.CategoryFilter) (*entity.Category, error) {
	return nil, nil
}

func (uc *CategoryUseCase) GetCategories(ctx context.Context, f *entity.CategoryFilter) ([]*entity.Category, error) {
	return nil, nil
}
