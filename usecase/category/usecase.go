package category

import (
	"context"
	"time"

	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/pkg/errutil"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/rs/zerolog/log"
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

	id, err := uc.categoryRepo.Create(ctx, c)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save new category to repo, err: %v", err)
		return err
	}

	c.CatID = goutil.String(id)

	return nil
}

func (uc *CategoryUseCase) UpdateCategory(ctx context.Context, c *entity.Category) error {
	now := time.Now().Unix()

	c.UpdateTime = goutil.Uint64(uint64(now))

	return nil
}

func (uc *CategoryUseCase) GetCategory(ctx context.Context, f *entity.CategoryFilter) (*entity.Category, error) {
	ec, err := uc.categoryRepo.Get(ctx, f)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get category from repo, err: %v", err)
		if err == errutil.ErrNotFound {
			return nil, errutil.NotFoundError(err)
		}
		return nil, err
	}

	return ec, nil
}

func (uc *CategoryUseCase) GetCategories(ctx context.Context, f *entity.CategoryFilter) ([]*entity.Category, error) {
	ecs, err := uc.categoryRepo.GetMany(ctx, f)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get categories from repo, err: %v", err)
		return nil, err
	}

	return ecs, nil
}
