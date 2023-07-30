package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var CreateBudgetValidator = validator.MustForm(map[string]validator.Validator{
	"category_id": &validator.String{
		Optional: false,
	},
	"amount": &validator.String{
		Optional:   false,
		Validators: []validator.StringFunc{entity.CheckPositiveMonetaryStr},
	},
	"budget_type": &validator.UInt32{
		Optional:   false,
		Validators: []validator.UInt32Func{entity.CheckBudgetType},
	},
	"budget_date": &validator.String{
		Optional:   true,
		Validators: []validator.StringFunc{entity.CheckDateStr},
	},
	"budget_repeat": &validator.UInt32{
		Optional:   true,
		Validators: []validator.UInt32Func{entity.CheckBudgetRepeat},
	},
})

func (h *budgetHandler) CreateBudget(ctx context.Context, req *presenter.CreateBudgetRequest, res *presenter.CreateBudgetResponse) error {
	userID := util.GetUserIDFromCtx(ctx)

	useCaseRes, err := h.budgetUseCase.CreateBudget(ctx, req.ToUseCaseReq(userID))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to create budget, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
