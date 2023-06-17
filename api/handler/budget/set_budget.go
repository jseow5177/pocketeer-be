package budget

import (
	"context"
	"fmt"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
)

var SetBudgetValidator = validator.MustForm(map[string]validator.Validator{
	"budget_name": &validator.String{
		Optional: false,
	},
	"budget_type": &validator.UInt32{
		Optional:   false,
		Validators: []validator.UInt32Func{entity.CheckBudgetType},
	},
	"budget_amount": &validator.String{
		Optional:   false,
		Validators: []validator.StringFunc{entity.CheckMonetaryStr},
	},
	"category_ids": &validator.Slice{
		MinLen: 1,
	},
	"range_start_date": &validator.String{
		Optional:   false,
		Validators: []validator.StringFunc{util.ValidateDateStr},
	},
	"range_end_date": &validator.String{
		Optional:   false,
		Validators: []validator.StringFunc{util.ValidateDateStr},
	},
})

func (h *budgetHandler) SetBudget(
	ctx context.Context,
	req *presenter.SetBudgetRequest,
	res *presenter.SetBudgetResponse,
) error {
	userID := util.GetUserIDFromCtx(ctx)
	if err := validateSetBudgetReq(req); err != nil {
		return err
	}

	useCaseRes, err := h.budgetUseCase.SetBudget(
		ctx,
		req.ToUseCaseReq(userID),
	)
	if err != nil {
		return err
	}

	res.Set(useCaseRes)
	return nil
}

func validateSetBudgetReq(
	req *presenter.SetBudgetRequest,
) error {
	if req.GetRangeEndDate() < req.GetRangeStartDate() {
		return fmt.Errorf("end date cannot be smaller than start date")
	}

	return nil
}
