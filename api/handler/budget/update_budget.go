package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var UpdateBudgetValidator = validator.MustForm(map[string]validator.Validator{
	"category_id": &validator.String{
		Optional: false,
	},
	"budget_date": &validator.String{
		Optional:   false,
		Validators: []validator.StringFunc{entity.CheckDateStr},
	},
	"amount": &validator.String{
		Optional:   true,
		Validators: []validator.StringFunc{entity.CheckPositiveMonetaryStr},
	},
	"budget_type": &validator.UInt32{
		Optional:   true,
		Validators: []validator.UInt32Func{entity.CheckBudgetType},
	},
	"budget_repeat": &validator.UInt32{
		Optional:   true,
		Validators: []validator.UInt32Func{entity.CheckBudgetRepeat},
	},
})

func (h *budgetHandler) UpdateBudget(ctx context.Context, req *presenter.UpdateBudgetRequest, res *presenter.UpdateBudgetResponse) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.budgetUseCase.UpdateBudget(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to update budget, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
