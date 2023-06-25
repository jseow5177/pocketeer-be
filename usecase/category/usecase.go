package category

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/rs/zerolog/log"
)

type categoryUseCase struct {
	categoryRepo repo.CategoryRepo
}

func NewCategoryUseCase(categoryRepo repo.CategoryRepo) UseCase {
	return &categoryUseCase{
		categoryRepo,
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

	_, err := uc.categoryRepo.Create(ctx, c)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save new category to repo, err: %v", err)
		return nil, err
	}

	return &CreateCategoryResponse{
		Category: c,
	}, nil
}

func (uc *categoryUseCase) UpdateCategory(ctx context.Context, req *UpdateCategoryRequest) (*UpdateCategoryResponse, error) {
	c, err := uc.categoryRepo.Get(ctx, req.ToCategoryFilter())
	if err != nil {
		return nil, err
	}

	nc := c.GetUpdates(req.ToCategoryUpdate(), true)
	if nc == nil {
		log.Ctx(ctx).Info().Msg("category has no updates")
		return &UpdateCategoryResponse{
			c,
		}, nil
	}

	if err = uc.categoryRepo.Update(ctx, req.ToCategoryFilter(), nc); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save category updates to repo, err: %v", err)
		return nil, err
	}

	return &UpdateCategoryResponse{
		c,
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
