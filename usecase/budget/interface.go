package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/data/presenter"
)

type UseCase interface {
	GetCategoryBudgetsByMonth(
		ctx context.Context,
		userID string,
		req *presenter.GetCategoryBudgetsByMonthRequest,
	) ([]*entity.Budget, error)

	GetBudgetBreakdownByYear(
		ctx context.Context,
		userID string,
		req *presenter.GetBudgetBreakdownByYearRequest,
	) (budgetBreakdown *entity.YearBudgetBreakdown, err error)

	SetBudget(
		ctx context.Context,
		userID string,
		req *presenter.SetBudgetRequest,
	) error
}
