package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

func NewCreateBudgetValidator(optionalCategoryID bool) validator.Validator {
	return validator.OptionalForm(map[string]validator.Validator{
		"category_id": &validator.String{
			Optional: optionalCategoryID,
		},
		"budget_date": &validator.String{
			Optional:   false,
			Validators: []validator.StringFunc{entity.CheckDateStr},
		},
		"amount": &validator.String{
			Optional:   false,
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
}

func (h *budgetHandler) CreateBudget(ctx context.Context, req *presenter.CreateBudgetRequest, res *presenter.CreateBudgetResponse) error {
	user := entity.GetUserFromCtx(ctx)
	req.Currency = user.Meta.Currency // TODO: Support currency on budget creation

	useCaseRes, err := h.budgetUseCase.CreateBudget(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to create budget, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
