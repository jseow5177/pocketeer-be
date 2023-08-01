package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var GetBudgetsValidator = validator.MustForm(map[string]validator.Validator{
	"category_ids": &validator.Slice{
		Optional:  false,
		Validator: &validator.String{},
	},
	"budget_date": &validator.String{
		Optional:   false,
		Validators: []validator.StringFunc{entity.CheckDateStr},
	},
	"timezone": &validator.String{
		Optional:   false,
		Validators: []validator.StringFunc{entity.CheckTimezone},
	},
})

func (h *budgetHandler) GetBudgets(ctx context.Context, req *presenter.GetBudgetsRequest, res *presenter.GetBudgetsResponse) error {
	userID := util.GetUserIDFromCtx(ctx)

	useCaseRes, err := h.budgetUseCase.GetBudgets(ctx, req.ToUseCaseReq(userID))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get budget, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
