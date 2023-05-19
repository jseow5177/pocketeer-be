package budget

import (
	"context"
	"fmt"

	"github.com/jseow5177/pockteer-be/api/middleware"
	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/errutil"
)

func (h *budgetHandler) GetAnnualBudgetBreakdown(
	ctx context.Context,
	req *presenter.GetAnnualBudgetBreakdownRequest,
	res *presenter.GetAnnualBudgetBreakdownResponse,
) error {
	userID := middleware.GetUserIDFromCtx(ctx)
	if userID == "" {
		return errutil.BadRequestError(
			fmt.Errorf("userID is not passed in ctx"),
		)
	}

	usecaseRes, err := h.budgetUseCase.GetAnnualBudgetBreakdown(
		ctx,
		req.ToUseCaseRes(userID),
	)
	if err != nil {
		return err
	}

	res.Set(usecaseRes)
	return nil
}
