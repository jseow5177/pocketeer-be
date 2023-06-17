package aggr

import (
	"context"
	"fmt"

	"github.com/jseow5177/pockteer-be/usecase/budget"
	"github.com/jseow5177/pockteer-be/usecase/category"
	"github.com/rs/zerolog/log"
)

type aggrUseCase struct {
	budgetUseCase   budget.UseCase
	categoryUseCase category.UseCase
}

func NewAggrUseCase(
	budgetUseCase budget.UseCase,
	categoryUseCase category.UseCase,
) UseCase {
	return &aggrUseCase{
		budgetUseCase,
		categoryUseCase,
	}
}

func (uc *aggrUseCase) GetBudgetWithCategories(
	ctx context.Context,
	req *GetBudgetWithCategoriesRequest,
) (*GetBudgetWithCategoriesResponse, error) {
	budgetRes, err := uc.budgetUseCase.GetBudget(
		ctx,
		&budget.GetBudgetRequest{
			UserID:   req.UserID,
			BudgetID: req.BudgetID,
			Date:     req.Date,
		},
	)
	if err != nil {
		return nil, err
	}

	budget := budgetRes.GetBudget()
	if budget == nil {
		log.Ctx(ctx).Error().Msgf("cannot find budget with budgetID=%s", req.GetBudgetID())
		return nil, fmt.Errorf("cannot find budget with budgetID=%s", req.GetBudgetID())
	}

	categoryRes, err := uc.categoryUseCase.GetCategories(
		ctx,
		&category.GetCategoriesRequest{
			CategoryIDs: budget.GetCategoryIDs(),
		},
	)
	if err != nil {
		return nil, err
	}
	if len(categoryRes.GetCategories()) != len(budget.GetCategoryIDs()) {
		log.Ctx(ctx).Error().Msgf("some categories are missing for ids=%+v", budget.GetCategoryIDs())
		return nil, fmt.Errorf("some categories are missing for ids=%+v", budget.GetCategoryIDs())
	}

	return &GetBudgetWithCategoriesResponse{
		Budget:     budget,
		Categories: categoryRes.GetCategories(),
	}, nil
}
