package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/data/presenter"
	"github.com/jseow5177/pockteer-be/dep/repo"
)

type budgetUseCase struct {
	budgetRepo repo.BudgetRepo
	catRepo    repo.CategoryRepo
}

func NewBudgetUseCase(
	budgetRepo repo.BudgetRepo,
	catRepo repo.CategoryRepo,
) UseCase {
	return &budgetUseCase{
		budgetRepo: budgetRepo,
		catRepo:    catRepo,
	}
}

func (uc *budgetUseCase) GetCategoryBudgetsByMonth(
	ctx context.Context,
	userID string,
	req *presenter.GetCategoryBudgetsByMonthRequest,
) ([]*entity.Budget, error) {
	categories, err := uc.catRepo.GetMany(ctx, req.ToCategoryFilter(userID))
	if err != nil {
		return nil, err
	}

	monthBudgets, err := uc.budgetRepo.GetMany(ctx, req.ToBudgetFilter(userID))
	if err != nil {
		return nil, err
	}

	categoryIDToBudget := getCategoryIDToBudgetMap(monthBudgets)
	for _, category := range categories {
		_, budgetSet := categoryIDToBudget[category.GetCategoryID()]
		if !budgetSet {
			categoryIDToBudget[category.GetCategoryID()] = nil
		}
	}

	return nil, nil
}

func (uc *budgetUseCase) GetBudgetBreakdownByYear(
	ctx context.Context,
	userID string,
	req *presenter.GetBudgetBreakdownByYearRequest,
) (budgetBreakdown *entity.YearBudgetBreakdown, err error) {
	fullYearBudgets, err := uc.budgetRepo.GetMany(ctx, nil)
	if err != nil {
		return nil, err
	}

	budgetNotSet := len(fullYearBudgets) == 0
	if budgetNotSet {
		return nil, nil
	}

	budgetBreakdown = entity.NewYearBudgetBreakdown(fullYearBudgets)

	return budgetBreakdown, nil
}

func (uc *budgetUseCase) SetBudget(
	ctx context.Context,
	userID string,
	req *presenter.SetBudgetRequest,
) error {
	fullYearBudgets, err := uc.budgetRepo.GetMany(ctx, nil)
	if err != nil {
		return err
	}

	var budgetBreakdown *entity.YearBudgetBreakdown
	if len(fullYearBudgets) == 0 {
		budgetBreakdown = entity.DefaultYearBudgetBreakdown(
			userID,
			req.GetCategoryID(),
			req.GetYear(),
		)
	} else {
		budgetBreakdown = entity.NewYearBudgetBreakdown(fullYearBudgets)
	}

	if req.BudgetType != nil {
		budgetBreakdown.SetBudgetType(req.GetBudgetType())
	}

	if req.BudgetAmount != nil {
		if req.IsDefault != nil {
			budgetBreakdown.SetDefaultBudget(req.GetBudgetAmount())
		} else {
			budgetBreakdown.SetMonthlyBudget(req.GetBudgetAmount())
		}
	}

	err = uc.budgetRepo.Set(ctx, budgetBreakdown.ToBudgets())
	if err != nil {
		return err
	}

	return nil
}

func getCategoryIDToBudgetMap(
	budgets []*entity.Budget,
) map[string]*entity.Budget {
	_map := make(map[string]*entity.Budget)
	for _, budget := range budgets {
		_map[budget.GetCategoryID()] = budget
	}

	return _map
}
