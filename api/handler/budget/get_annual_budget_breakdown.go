package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
)

var GetAnnualBudgetBreakdownValidator = validator.MustForm(map[string]validator.Validator{
	"year": &validator.UInt32{
		Optional: false,
	},
	"category_id": &validator.String{
		Optional: false,
	},
})

func (h *budgetHandler) GetAnnualBudgetBreakdown(
	ctx context.Context,
	req *presenter.GetAnnualBudgetBreakdownRequest,
	res *presenter.GetAnnualBudgetBreakdownResponse,
) error {
	userID := util.GetUserIDFromCtx(ctx)

	usecaseRes, err := h.budgetUseCase.GetAnnualBudgetBreakdown(
		ctx,
		req.ToUseCaseReq(userID),
	)
	if err != nil {
		return err
	}

	res.Set(usecaseRes)
	return nil
}
