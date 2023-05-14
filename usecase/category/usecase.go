package category

import (
	"context"
	"errors"
	"time"

	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/data/presenter"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/pkg/errutil"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/rs/zerolog/log"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

type CategoryUseCase struct {
	categoryRepo repo.CategoryRepo
}

func NewCategoryUseCase(categoryRepo repo.CategoryRepo) UseCase {
	return &CategoryUseCase{
		categoryRepo: categoryRepo,
	}
}

func (uc *CategoryUseCase) CreateCategory(ctx context.Context, userID string, req *presenter.CreateCategoryRequest) (*entity.Category, error) {
	var (
		c   = req.ToCategoryEntity(userID)
		now = uint64(time.Now().Unix())
	)

	c.CreateTime = goutil.Uint64(now)
	c.UpdateTime = goutil.Uint64(now)

	id, err := uc.categoryRepo.Create(ctx, c)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save new category to repo, err: %v", err)
		return nil, err
	}

	c.CategoryID = goutil.String(id)

	return c, nil
}

func (uc *CategoryUseCase) UpdateCategory(ctx context.Context, userID string, req *presenter.UpdateCategoryRequest) (*entity.Category, error) {
	c, err := uc.GetCategory(ctx, userID, req.ToGetCategoryRequest())
	if err != nil {
		return nil, err
	}

	nc := uc.getCategoryUpdates(c, req.ToCategoryEntity())
	if nc == nil {
		// no updates
		log.Ctx(ctx).Info().Msg("category has no updates")
		return c, nil
	}

	cf := req.ToCategoryFilter(userID)
	if err = uc.categoryRepo.Update(ctx, cf, nc); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save category updates to repo, err: %v", err)
		return nil, err
	}

	// merge
	goutil.MergeWithPtrFields(c, nc)

	return c, nil
}

func (uc *CategoryUseCase) getCategoryUpdates(old, changes *entity.Category) *entity.Category {
	var hasUpdates bool

	nc := &entity.Category{
		UpdateTime: goutil.Uint64(uint64(time.Now().Unix())),
	}

	if changes.CategoryName != nil && changes.GetCategoryName() != old.GetCategoryName() {
		hasUpdates = true
		nc.CategoryName = changes.CategoryName
	}

	if !hasUpdates {
		return nil
	}

	return nc
}

func (uc *CategoryUseCase) GetCategory(ctx context.Context, userID string, req *presenter.GetCategoryRequest) (*entity.Category, error) {
	cf := req.ToCategoryFilter(userID)

	c, err := uc.categoryRepo.Get(ctx, cf)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get category from repo, err: %v", err)
		if err == errutil.ErrNotFound {
			return nil, errutil.NotFoundError(ErrCategoryNotFound)
		}
		return nil, err
	}

	return c, nil
}

func (uc *CategoryUseCase) GetCategories(ctx context.Context, userID string, req *presenter.GetCategoriesRequest) ([]*entity.Category, error) {
	cf := req.ToCategoryFilter(userID)

	cs, err := uc.categoryRepo.GetMany(ctx, cf)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get categories from repo, err: %v", err)
		return nil, err
	}

	return cs, nil
}
