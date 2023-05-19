package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/middleware"
	"github.com/jseow5177/pockteer-be/api/presenter"
)

func (h *budgetHandler) GetCategoryBudgetsByMonth(
	ctx context.Context,
	req *presenter.GetCategoryBudgetsByMonthRequest,
	res *presenter.GetCategoryBudgetsByMonthResponse,
) error {
	var (
		userID = middleware.GetUserIDFromCtx(ctx)
	)

	usecaseRes, err := h.budgetUseCase.GetCategoryBudgetsByMonth(
		ctx,
		userID,
		req,
	)
	if err != nil {
		return err
	}

	res.Budgets = toBudgets(usecaseRes.CategoryBudgets)
	return nil
}
