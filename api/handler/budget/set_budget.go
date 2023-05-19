package budget

import (
	"context"
	"fmt"

	"github.com/jseow5177/pockteer-be/api/middleware"
	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/errutil"
	"github.com/jseow5177/pockteer-be/pkg/validator"
)

var SetBudgetValidator = validator.MustForm(map[string]validator.Validator{
	"category_id": &validator.String{
		Optional: false,
	},
	"year": &validator.UInt32{
		Optional: false,
	},
})

func (h *budgetHandler) SetBudget(
	ctx context.Context,
	req *presenter.SetBudgetRequest,
	res *presenter.SetBudgetResponse,
) error {
	userID := middleware.GetUserIDFromCtx(ctx)
	if userID == "" {
		return errutil.BadRequestError(
			fmt.Errorf("userID is not passed in ctx"),
		)
	}

	usecaseRes, err := h.budgetUseCase.SetBudget(
		ctx,
		req.ToUseCaseReq(userID),
	)
	if err != nil {
		return err
	}

	res.Set(usecaseRes)
	return nil
}
