package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/data/presenter"
)

type UseCase interface {
	GetMonthBudgets(
		ctx context.Context,
		userID string,
		req *presenter.GetMonthBudgetsRequest,
	) ([]*entity.Budget, error)

	GetFullYearBudget(
		ctx context.Context,
		userID string,
		req *presenter.GetFullYearBudgetRequest,
	) (defaultBudget *entity.Budget, monthlyBudgets []*entity.Budget, err error)

	SetBudget(
		ctx context.Context,
		userID string,
		req *presenter.SetBudgetRequest,
	) error
}
