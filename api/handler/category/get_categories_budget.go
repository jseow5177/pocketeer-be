package category

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var GetCategoriesBudgetValidator = validator.MustForm(map[string]validator.Validator{
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

func (h *categoryHandler) GetCategoriesBudget(
	ctx context.Context,
	req *presenter.GetCategoriesBudgetRequest,
	res *presenter.GetCategoriesBudgetResponse,
) error {
	userID := util.GetUserIDFromCtx(ctx)

	useCaseRes, err := h.categoryUseCase.GetCategoriesBudget(ctx, req.ToUseCaseReq(userID))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get categories budget, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
