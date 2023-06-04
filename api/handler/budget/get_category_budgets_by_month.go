package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
)

var GetCategoryBudgetsByMonthValidator = validator.MustForm(map[string]validator.Validator{
	"year": &validator.UInt32{
		Optional: false,
	},
	"month": &validator.UInt32{
		Optional: false,
	},
})

func (h *budgetHandler) GetCategoryBudgetsByMonth(
	ctx context.Context,
	req *presenter.GetCategoryBudgetsByMonthRequest,
	res *presenter.GetCategoryBudgetsByMonthResponse,
) error {
	userID := util.GetUserIDFromCtx(ctx)

	usecaseRes, err := h.budgetUseCase.GetCategoryBudgetsByMonth(
		ctx,
		req.ToUseCaseReq(userID),
	)
	if err != nil {
		return err
	}

	res.Set(usecaseRes)
	return nil
}
