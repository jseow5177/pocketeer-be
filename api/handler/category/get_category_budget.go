package category

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var GetCategoryBudgetValidator = validator.MustForm(map[string]validator.Validator{
	"category_id": &validator.String{
		Optional: false,
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

func (h *categoryHandler) GetCategoryBudget(
	ctx context.Context,
	req *presenter.GetCategoryBudgetRequest,
	res *presenter.GetCategoryBudgetResponse,
) error {
	userID := util.GetUserIDFromCtx(ctx)

	useCaseRes, err := h.categoryUseCase.GetCategoryBudget(ctx, req.ToUseCaseReq(userID))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get category budget, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
