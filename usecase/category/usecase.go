package category

import (
	"context"
	"fmt"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/usecase/budget"
	"github.com/rs/zerolog/log"
)

type categoryUseCase struct {
	categoryRepo    repo.CategoryRepo
	transactionRepo repo.TransactionRepo
	budgetUseCase   budget.UseCase
}

func NewCategoryUseCase(categoryRepo repo.CategoryRepo, transactionRepo repo.TransactionRepo, budgetUseCase budget.UseCase) UseCase {
	return &categoryUseCase{
		categoryRepo,
		transactionRepo,
		budgetUseCase,
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

func (uc *categoryUseCase) GetCategoryBudget(ctx context.Context, req *GetCategoryBudgetRequest) (*GetCategoryBudgetResponse, error) {
	c, err := uc.categoryRepo.Get(ctx, req.ToCategoryFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get category from repo, err: %v", err)
		return nil, err
	}

	if !c.CanAddBudget() {
		return &GetCategoryBudgetResponse{
			Category: c,
		}, nil
	}

	res, err := uc.budgetUseCase.GetBudget(ctx, req.ToGetBudgetRequest())
	if err != nil {
		return nil, err
	}
	c.SetBudget(res.Budget)

	return &GetCategoryBudgetResponse{
		Category: c,
	}, nil
}

func (uc *categoryUseCase) CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*CreateCategoryResponse, error) {
	_, err := uc.categoryRepo.Get(ctx, req.ToCategoryFilter())
	if err != nil && err != repo.ErrCategoryNotFound {
		log.Ctx(ctx).Error().Msgf("fail to get category from repo, err: %v", err)
		return nil, err
	}

	if err == nil {
		return nil, repo.ErrCategoryAlreadyExists
	}

	c, err := req.ToCategoryEntity()
	if err != nil {
		return nil, err
	}

	_, err = uc.categoryRepo.Create(ctx, c)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save new category to repo, err: %v", err)
		return nil, err
	}

	return &CreateCategoryResponse{
		Category: c,
	}, nil
}

func (uc *categoryUseCase) CreateCategories(ctx context.Context, req *CreateCategoriesRequest) (*CreateCategoriesResponse, error) {
	var (
		cs  = make([]*entity.Category, 0)
		uks = make(map[string]bool)
	)

	for _, r := range req.Categories {
		if _, ok := uks[fmt.Sprintf("%v-%v", r.GetCategoryName(), r.GetCategoryType())]; ok {
			return nil, repo.ErrCategoryAlreadyExists
		}

		_, err := uc.categoryRepo.Get(ctx, r.ToCategoryFilter())
		if err != nil && err != repo.ErrCategoryNotFound {
			log.Ctx(ctx).Error().Msgf("fail to get category from repo, err: %v", err)
			return nil, err
		}

		if err == nil {
			return nil, repo.ErrCategoryAlreadyExists
		}

		c, err := r.ToCategoryEntity()
		if err != nil {
			return nil, err
		}

		cs = append(cs, c)
		uks[c.GetCategoryName()] = true
	}

	if _, err := uc.categoryRepo.CreateMany(ctx, cs); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save new categories to repo, err: %v", err)
		return nil, err
	}

	return &CreateCategoriesResponse{
		Categories: cs,
	}, nil
}

func (uc *categoryUseCase) UpdateCategory(ctx context.Context, req *UpdateCategoryRequest) (*UpdateCategoryResponse, error) {
	c, err := uc.categoryRepo.Get(ctx, req.ToCategoryFilter())
	if err != nil {
		return nil, err
	}

	cu, hasUpdate, err := c.Update(req.ToCategoryUpdate())
	if err != nil {
		return nil, err
	}

	if !hasUpdate {
		log.Ctx(ctx).Info().Msg("category has no updates")
		return &UpdateCategoryResponse{
			c,
		}, nil
	}

	if err = uc.categoryRepo.Update(ctx, req.ToCategoryFilter(), cu); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save category updates to repo, err: %v", err)
		return nil, err
	}

	return &UpdateCategoryResponse{
		Category: c,
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

func (uc *categoryUseCase) GetCategoriesBudget(ctx context.Context, req *GetCategoriesBudgetRequest) (*GetCategoriesBudgetResponse, error) {
	cs, err := uc.categoryRepo.GetMany(ctx, req.ToCategoryFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get categories from repo, err: %v", err)
		return nil, err
	}

	var (
		catIDs = make([]string, 0)
		cbs    = make(map[string]*entity.Category)
	)
	for _, c := range cs {
		if c.CanAddBudget() {
			catIDs = append(catIDs, c.GetCategoryID())
			cbs[c.GetCategoryID()] = c
		}
	}

	if len(catIDs) == 0 {
		return &GetCategoriesBudgetResponse{
			Categories: cs,
		}, nil
	}

	res, err := uc.budgetUseCase.GetBudgets(ctx, req.ToGetBudgetsRequest())
	if err != nil {
		return nil, err
	}

	for _, b := range res.Budgets {
		if c, ok := cbs[b.GetCategoryID()]; ok {
			c.SetBudget(b)
		}
	}

	return &GetCategoriesBudgetResponse{
		Categories: cs,
	}, nil
}
