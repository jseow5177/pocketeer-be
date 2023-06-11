package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
)

var GetBudgetValidator = validator.MustForm(map[string]validator.Validator{
	"budget_id": &validator.String{
		Optional: false,
	},
	"date": &validator.String{
		Optional:   false,
		Validators: []validator.StringFunc{util.ValidateDateStr},
	},
})

func (h *budgetHandler) GetBudget(
	ctx context.Context,
	req *presenter.GetBudgetRequest,
	res *presenter.GetBudgetResponse,
) error {
	userID := util.GetUserIDFromCtx(ctx)
	useCaseRes, err := h.aggrUsecase.GetBudgetWithCategories(
		ctx,
		req.ToUseCaseReq(userID),
	)
	if err != nil {
		return err
	}

	res.Set(useCaseRes)
	return nil
}
