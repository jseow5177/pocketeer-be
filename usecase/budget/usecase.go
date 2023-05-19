package budget

import (
	"context"

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
	req *GetCategoryBudgetsByMonthRequest,
) (*GetCategoryBudgetsByMonthResponse, error) {
	categories, err := uc.categoryRepo.GetMany(ctx, req.ToCategoryFilter())
	if err != nil {
		return nil, err
	}

	monthBudgets, err := uc.budgetRepo.GetMany(ctx, req.ToBudgetFilter())
	if err != nil {
		return nil, err
	}

	categoryBudgets, err := uc.mergeCategoryAndBudget(categories, monthBudgets)
	if err != nil {
		return nil, err
	}

	return &GetCategoryBudgetsByMonthResponse{
		CategoryBudgets: categoryBudgets,
	}, nil
}

func (uc *budgetUseCase) mergeCategoryAndBudget(
	categories []*entity.Category,
	monthBudgets []*entity.Budget,
) ([]*CategoryBudget, error) {
	catIDToBudget := util.GetCategoryIDToBudgetMap(monthBudgets)
	catIDToCategoryBudget := make(map[string]*CategoryBudget)

	for _, category := range categories {
		budget, _ := catIDToBudget[category.GetCategoryID()]

		catIDToCategoryBudget[category.GetCategoryID()] = &CategoryBudget{
			Category: category,
			Budget:   budget,
		}
	}

	categoryBudgets := make([]*CategoryBudget, 0)
	for _, categoryBudget := range catIDToCategoryBudget {
		categoryBudgets = append(categoryBudgets, categoryBudget)
	}

	return categoryBudgets, nil
}

func (uc *budgetUseCase) GetAnnualBudgetBreakdown(
	ctx context.Context,
	req *GetAnnualBudgetBreakdownRequest,
) (*GetAnnualBudgetBreakdownResponse, error) {
	fullYearBudgets, err := uc.budgetRepo.GetMany(ctx, req.ToFullBudgetFilter())
	if err != nil {
		return nil, err
	}

	var (
		budgetNotSet    = len(fullYearBudgets) == 0
		budgetBreakdown *entity.AnnualBudgetBreakdown
	)

	if budgetNotSet {
		budgetBreakdown = entity.DefaultAnnualBudgetBreakdown(
			req.GetUserID(),
			req.GetCategoryID(),
			req.GetYear(),
		)
	} else {
		budgetBreakdown, err = entity.NewAnnualBudgetBreakdown(fullYearBudgets)
		if err != nil {
			return nil, err
		}
	}

	return &GetAnnualBudgetBreakdownResponse{
		AnnualBudgetBreakdown: budgetBreakdown,
	}, nil
}

func (uc *budgetUseCase) SetBudget(
	ctx context.Context,
	req *SetBudgetRequest,
) (*SetBudgetResponse, error) {
	fullYearBudgets, err := uc.budgetRepo.GetMany(ctx, req.ToFullBudgetFilter())
	if err != nil {
		return nil, err
	}

	var (
		budgetNotSet    = len(fullYearBudgets) == 0
		budgetBreakdown *entity.AnnualBudgetBreakdown
	)

	if budgetNotSet {
		budgetBreakdown = entity.DefaultAnnualBudgetBreakdown(
			req.GetUserID(),
			req.GetCategoryID(),
			req.GetYear(),
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
