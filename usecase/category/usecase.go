package category

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/rs/zerolog/log"
)

type categoryUseCase struct {
	categoryRepo repo.CategoryRepo
}

func NewCategoryUseCase(categoryRepo repo.CategoryRepo) UseCase {
	return &categoryUseCase{
		categoryRepo: categoryRepo,
	}
}

func (uc *categoryUseCase) GetCategory(ctx context.Context, req *GetCategoryRequest) (*GetCategoryResponse, error) {
	c, err := uc.categoryRepo.Get(ctx, req.ToCategoryFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get category from repo, err: %v", err)
		return nil, err
	}

	return &GetCategoryResponse{
		Category: c,
	}, nil
}

func (uc *categoryUseCase) CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*CreateCategoryResponse, error) {
	c := req.ToCategoryEntity()

	id, err := uc.categoryRepo.Create(ctx, c)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save new category to repo, err: %v", err)
		return nil, err
	}

	c.CategoryID = goutil.String(id)

	return &CreateCategoryResponse{
		Category: c,
	}, nil
}

func (uc *categoryUseCase) UpdateCategory(ctx context.Context, req *UpdateCategoryRequest) (*UpdateCategoryResponse, error) {
	c, err := uc.GetCategory(ctx, req.ToGetCategoryRequest())
	if err != nil {
		return nil, err
	}
	category := c.Category

	nc := uc.getCategoryUpdates(category, req.ToCategoryEntity())
	if nc == nil {
		// no updates
		log.Ctx(ctx).Info().Msg("category has no updates")
		return &UpdateCategoryResponse{
			Category: category,
		}, nil
	}

	if err = uc.categoryRepo.Update(ctx, req.ToCategoryFilter(), nc); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save category updates to repo, err: %v", err)
		return nil, err
	}

	// merge
	goutil.MergeWithPtrFields(category, nc)

	return &UpdateCategoryResponse{
		Category: category,
	}, nil
}

func (uc *categoryUseCase) GetCategories(ctx context.Context, req *GetCategoriesRequest) (*GetCategoriesResponse, error) {
	cs, err := uc.categoryRepo.GetMany(ctx, req.ToCategoryFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get categories from repo, err: %v", err)
		return nil, err
	}

	return &GetCategoriesResponse{
		Categories: cs,
	}, nil
}

func (uc *categoryUseCase) getCategoryUpdates(old, changes *entity.Category) *entity.Category {
	var hasUpdates bool

	nc := new(entity.Category)

	if changes.CategoryName != nil && changes.GetCategoryName() != old.GetCategoryName() {
		hasUpdates = true
		nc.CategoryName = changes.CategoryName
	}

	if !hasUpdates {
		return nil
	}

	return nc
}
