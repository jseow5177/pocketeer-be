package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/middleware"
	"github.com/jseow5177/pockteer-be/data/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/usecase/budget"
)

var SetBudgetValidator = validator.MustForm(map[string]validator.Validator{
	"cat_id": &validator.String{
		Optional: false,
	},
	"budget_type": &validator.UInt32{
		Optional: false,
	},
	"curr_year": &validator.UInt32{
		Optional: false,
	},
	"curr_month": &validator.UInt32{
		Optional: false,
	},
})

func (h *BudgetHandler) SetBudget(
	ctx context.Context,
	req *presenter.SetBudgetRequest,
	res *presenter.SetBudgetResponse,
) error {
	var (
		userID = middleware.GetUserIDFromCtx(ctx)
		uc     = budget.NewBudgetUseCase(h.budgetConfigRepo, h.categoryRepo)
	)

	err := uc.SetBudget(ctx, userID, req)
	if err != nil {
		return err
	}

	return nil
}
