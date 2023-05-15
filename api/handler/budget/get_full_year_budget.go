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

var GetFullYearBudgetValidator = validator.MustForm(map[string]validator.Validator{
	"year": &validator.UInt32{
		Optional: false,
	},
	"cat_id": &validator.String{
		Optional: false,
	},
})

func (h *BudgetHandler) GetFullYearBudget(
	ctx context.Context,
	req *presenter.GetFullYearBudgetRequest,
	res *presenter.GetFullYearBudgetResponse,
) error {
	var (
		userID = middleware.GetUserIDFromCtx(ctx)
		uc     = budget.NewBudgetUseCase(h.budgetConfigRepo, h.categoryRepo)
	)

	defaultBudget, monthlyBudgets, err := uc.GetFullYearBudget(
		ctx, userID, req,
	)
	if err != nil {
		return err
	}

	res.FullYearBudget = toPresenterFullYearBudget(defaultBudget, monthlyBudgets)
	return nil
}

func toPresenterFullYearBudget(
	defaultBudget *entity.Budget,
	monthlyBudgets []*entity.Budget,
) *presenter.FullYearBudget {
	return &presenter.FullYearBudget{
		CatID:           goutil.String(defaultBudget.GetCatID()),
		BudgetType:      goutil.Uint32(defaultBudget.GetBudgetType()),
		TransactionType: goutil.Uint32(defaultBudget.GetTransactionType()),
		Year:            goutil.Uint32(defaultBudget.GetYear()),
		DefaultBudget:   goutil.Int64(defaultBudget.GetBudgetAmount()),
		MonthlyBudgets:  toBasicBudgets(monthlyBudgets),
	}
}

func toBasicBudgets(
	monthlyBudgets []*entity.Budget,
) []*presenter.BasicBudget {
	basicBudgets := make([]*presenter.BasicBudget, len(monthlyBudgets))
	for idx, budget := range monthlyBudgets {
		basicBudgets[idx] = &presenter.BasicBudget{
			Year:   goutil.Uint32(budget.GetYear()),
			Month:  goutil.Uint32(budget.GetMonth()),
			Amount: goutil.Int64(budget.GetBudgetAmount()),
		}
	}

	return basicBudgets
}
