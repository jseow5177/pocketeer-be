package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var DeleteBudgetValidator = validator.MustForm(map[string]validator.Validator{
	"category_id": &validator.String{
		Optional: false,
	},
	"budget_date": &validator.String{
		Optional:   false,
		Validators: []validator.StringFunc{entity.CheckDateStr},
	},
	"budget_repeat": &validator.UInt32{
		Optional:   true,
		Validators: []validator.UInt32Func{entity.CheckBudgetRepeat},
	},
})

func (h *budgetHandler) DeleteBudget(ctx context.Context, req *presenter.DeleteBudgetRequest, res *presenter.DeleteBudgetResponse) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.budgetUseCase.DeleteBudget(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to delete budget, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
