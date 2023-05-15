package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/middleware"
	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/data/presenter"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/usecase/budget"
)

var GetMonthBudgetsValidator = validator.MustForm(map[string]validator.Validator{
	"year": &validator.Uint32{
		Optional: false,
	},
	"month": &validator.Uint32{
		Optional: false,
	},
})

func (h *BudgetHandler) GetMonthBudgets(
	ctx context.Context,
	req *presenter.GetMonthBudgetsRequest,
	res *presenter.GetMonthBudgetsResponse,
) error {
	var (
		userID = middleware.GetUserIDFromCtx(ctx)
		uc     = budget.NewBudgetUseCase(h.budgetConfigRepo, h.categoryRepo)
	)

	budgets, err := uc.GetMonthBudgets(ctx, userID, req)
	if err != nil {
		return err
	}

	res.Budgets = toPresenterBudgets(budgets)
	return nil
}

func toPresenterBudgets(entBudgets []*entity.Budget) []*presenter.Budget {
	budgets := make([]*presenter.Budget, len(entBudgets))
	for idx, budget := range entBudgets {
		budgets[idx] = &presenter.Budget{
			CatID:           goutil.String(budget.GetCatID()),
			CatName:         goutil.String(budget.GetCatName()),
			BudgetType:      goutil.Uint32(budget.GetBudgetType()),
			TransactionType: goutil.Uint32(budget.GetTransactionType()),
			Year:            goutil.Uint32(budget.GetYear()),
			Month:           goutil.Uint32(budget.GetMonth()),
			BudgetAmount:    goutil.Int64(budget.GetBudgetAmount()),
			Used:            goutil.Int64(budget.GetUsed()),
		}
	}

	return budgets
}
