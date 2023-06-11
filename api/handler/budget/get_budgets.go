package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
)

var GetBudgetsValidator = validator.MustForm(map[string]validator.Validator{
	"date": &validator.String{
		Optional:   false,
		Validators: []validator.StringFunc{util.ValidateDateStr},
	},
})

func (h *budgetHandler) GetBudgets(
	ctx context.Context,
	req *presenter.GetBudgetsRequest,
	res *presenter.GetBudgetsResponse,
) error {
	userID := util.GetUserIDFromCtx(ctx)

	useCaseRes, err := h.budgetUseCase.GetBudgets(
		ctx,
		req.ToUseCaseReq(userID),
	)
	if err != nil {
		return err
	}

	res.Set(useCaseRes)
	return nil
}
