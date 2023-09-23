package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var GetBudgetValidator = validator.MustForm(map[string]validator.Validator{
	"category_id": &validator.String{
		Optional: false,
	},
	"budget_date": &validator.String{
		Optional:   false,
		Validators: []validator.StringFunc{entity.CheckDateStr},
	},
})

func (h *budgetHandler) GetBudget(ctx context.Context, req *presenter.GetBudgetRequest, res *presenter.GetBudgetResponse) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.budgetUseCase.GetBudget(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get budget, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
