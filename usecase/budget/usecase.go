package budget

import (
	"context"
	"fmt"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/rs/zerolog/log"
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
		budgetRepo,
		categoryRepo,
	}
}

func (uc *budgetUseCase) GetBudget(
	ctx context.Context,
	req *GetBudgetRequest,
) (*GetBudgetResponse, error) {
	budget, err := uc.budgetRepo.Get(ctx, req.ToBudgetFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get budget from repo, err: %v", err)
		return nil, err
	}

	budget.FilterBreakdownByDate(req.GetDate())

	return &GetBudgetResponse{
		Budget: budget,
	}, nil
}

func (uc *budgetUseCase) GetBudgets(
	ctx context.Context,
	req *GetBudgetsRequest,
) (*GetBudgetsResponse, error) {
	budgets, err := uc.budgetRepo.GetMany(ctx, req.ToBudgetFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get budgets from repo, err: %v", err)
		return nil, err
	}

	filteredBudgets := make([]*entity.Budget, 0)
	for _, budget := range budgets {
		if budget.IsBreakdownAvailable(req.GetDate()) {
			filteredBudgets = append(filteredBudgets, budget)
		}
	}

	return &GetBudgetsResponse{
		Budgets: filteredBudgets,
	}, nil
}

func (uc *budgetUseCase) SetBudget(
	ctx context.Context,
	req *SetBudgetRequest,
) (*SetBudgetResponse, error) {
	var (
		budget *entity.Budget
		err    error
	)

	if req.BudgetID != nil {
		budget, err = uc.budgetRepo.Get(ctx, req.ToBudgetFilter())
		if err != nil {
			log.Ctx(ctx).Error().Msgf("setBudget err, fail to get budget from repo, err: %v", err)
			return nil, err
		}
	} else {
		budget = entity.NewBudget(req.GetUserID(), req.GetBudgetType())
	}

	if req.BudgetType != nil {
		err = budget.SetBudgetType(req.GetBudgetType())
		if err != nil {
			log.Ctx(ctx).Error().Msgf("setBudgetType err: %v", err)
			return nil, err
		}
	}

	if req.BudgetName != nil {
		budget.SetBudgetName(req.GetBudgetName())
	}

	if req.CategoryIDs != nil {
		budget.SetCategoryIDs(req.GetCategoryIDs())
	}

	if req.BudgetAmount != nil {
		budget.SetBudgetAmount(
			req.GetBudgetAmount(),
			req.GetRangeStartDate(),
			req.GetRangeEndDate(),
		)
	}

	err = uc.budgetRepo.Set(ctx, budget)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to set budget with repo, err: %v", err)
		return nil, err
	}

	return &SetBudgetResponse{}, nil
}

func (uc *budgetUseCase) GetBudgetWithCategories(
	ctx context.Context,
	req *GetBudgetWithCategoriesRequest,
) (*GetBudgetWithCategoriesResponse, error) {
	budgetRes, err := uc.GetBudget(
		ctx,
		req.ToGetBudgetRequest(),
	)
	if err != nil {
		return nil, err
	}

	budget := budgetRes.GetBudget()
	if budget == nil {
		log.Ctx(ctx).Error().Msgf("cannot find budget with budgetID=%s", req.GetBudgetID())
		return nil, fmt.Errorf("cannot find budget with budgetID=%s", req.GetBudgetID())
	}

	cs, err := uc.categoryRepo.GetMany(
		ctx,
		&repo.CategoryFilter{
			CategoryIDs: budget.GetCategoryIDs(),
		},
	)
	if err != nil {
		return nil, err
	}

	if len(cs) != len(budget.GetCategoryIDs()) {
		log.Ctx(ctx).Error().Msgf("some categories are missing for ids=%+v", budget.GetCategoryIDs())
		return nil, fmt.Errorf("some categories are missing for ids=%+v", budget.GetCategoryIDs())
	}

	return &GetBudgetWithCategoriesResponse{
		Budget:     budget,
		Categories: cs,
	}, nil
}
