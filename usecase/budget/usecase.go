package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/usecase/util"
)

type budgetUseCase struct {
	budgetRepo   repo.BudgetRepo
	categoryRepo repo.CategoryRepo
}

func NewBudgetUseCase(
	budgetRepo repo.BudgetRepo,
	categoryRepo repo.CategoryRepo,
) UseCase {
	return &budgetUseCase{
		budgetRepo:   budgetRepo,
		categoryRepo: categoryRepo,
	}
}

func (uc *budgetUseCase) GetCategoryBudgetsByMonth(
	ctx context.Context,
	userID string,
	req *presenter.GetCategoryBudgetsByMonthRequest,
) (*GetCategoryBudgetsByMonthResponse, error) {
	categories, err := uc.categoryRepo.GetMany(ctx, req.ToCategoryFilter(userID))
	if err != nil {
		return nil, err
	}

	monthBudgets, err := uc.budgetRepo.GetMany(ctx, req.ToBudgetFilter(userID))
	if err != nil {
		return nil, err
	}

	categoryBudgets, err := uc.mergeCategoryAndBudget(categories, monthBudgets)
	if err != nil {
		return nil, err
	}

	return &GetCategoryBudgetsByMonthResponse{CategoryBudgets: categoryBudgets}, nil
}

func (uc *budgetUseCase) mergeCategoryAndBudget(
	categories []*entity.Category,
	monthBudgets []*entity.Budget,
) ([]*CategoryBudget, error) {
	categoryIDToBudget := util.GetCategoryIDToBudgetMap(monthBudgets)

	for _, category := range categories {
		_, isBudgetSet := categoryIDToBudget[category.GetCategoryID()]
		if !isBudgetSet {
			categoryIDToBudget[category.GetCategoryID()] = nil
		}
	}

	return toCategoryBudgets(categoryIDToBudget, categories)
}

func (uc *budgetUseCase) GetAnnualBudgetBreakdown(
	ctx context.Context,
	userID string,
	req *presenter.GetAnnualBudgetBreakdownRequest,
) (*GetAnnualBudgetBreakdownResponse, error) {
	fullYearBudgets, err := uc.budgetRepo.GetMany(ctx, req.ToFullBudgetFilter(userID))
	if err != nil {
		return nil, err
	}

	budgetNotSet := len(fullYearBudgets) == 0
	if budgetNotSet {
		return nil, nil
	}

	budgetBreakdown, err := entity.NewAnnualBudgetBreakdown(fullYearBudgets)
	if err != nil {
		return nil, err
	}

	return &GetAnnualBudgetBreakdownResponse{AnnualBudgetBreakdown: budgetBreakdown}, nil
}

func (uc *budgetUseCase) SetBudget(
	ctx context.Context,
	userID string,
	req *presenter.SetBudgetRequest,
) (*SetBudgetResponse, error) {
	fullYearBudgets, err := uc.budgetRepo.GetMany(ctx, req.ToFullBudgetFilter(userID))
	if err != nil {
		return nil, err
	}

	var budgetBreakdown *entity.AnnualBudgetBreakdown
	budgetNotSet := len(fullYearBudgets) == 0

	if budgetNotSet {
		budgetBreakdown = entity.DefaultAnnualBudgetBreakdown(
			userID, req.GetCategoryID(), req.GetYear(),
		)
	} else {
		budgetBreakdown, err = entity.NewAnnualBudgetBreakdown(fullYearBudgets)
		if err != nil {
			return nil, err
		}
	}

	if req.BudgetType != nil {
		budgetBreakdown.SetBudgetType(req.GetBudgetType())
	}

	if req.BudgetAmount != nil {
		if req.IsDefault != nil {
			budgetBreakdown.SetDefaultBudget(req.GetBudgetAmount())
		} else {
			err = budgetBreakdown.SetMonthlyBudget(req.GetMonth(), req.GetBudgetAmount())
			if err != nil {
				return nil, err
			}
		}
	}

	err = uc.budgetRepo.Set(ctx, budgetBreakdown.ToBudgets())
	if err != nil {
		return nil, err
	}

	return &SetBudgetResponse{AnnualBudgetBreakdown: budgetBreakdown}, nil
}
