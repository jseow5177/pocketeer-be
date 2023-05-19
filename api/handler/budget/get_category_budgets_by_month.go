package budget

import (
	"context"
	"fmt"

	"github.com/jseow5177/pockteer-be/api/middleware"
	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/errutil"
	"github.com/jseow5177/pockteer-be/pkg/validator"
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
	userID := middleware.GetUserIDFromCtx(ctx)
	if userID == "" {
		return errutil.BadRequestError(
			fmt.Errorf("userID is not passed in ctx"),
		)
	}

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
