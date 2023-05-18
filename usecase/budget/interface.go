package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
)

type UseCase interface {
	GetCategoryBudgetsByMonth(
		ctx context.Context,
		userID string,
		req *presenter.GetCategoryBudgetsByMonthRequest,
	) (*GetCategoryBudgetsByMonthResponse, error)

	GetAnnualBudgetBreakdown(
		ctx context.Context,
		userID string,
		req *presenter.GetAnnualBudgetBreakdownRequest,
	) (*GetAnnualBudgetBreakdownResponse, error)

	SetBudget(
		ctx context.Context,
		userID string,
		req *presenter.SetBudgetRequest,
	) (*SetBudgetResponse, error)
}

type CategoryBudget struct {
	CategoryName string
	Budget       *entity.Budget
}

type GetCategoryBudgetsByMonthResponse struct {
	CategoryBudgets []*CategoryBudget
}

type GetAnnualBudgetBreakdownResponse struct {
	AnnualBudgetBreakdown *entity.AnnualBudgetBreakdown
}

type SetBudgetResponse struct {
	AnnualBudgetBreakdown *entity.AnnualBudgetBreakdown
}
