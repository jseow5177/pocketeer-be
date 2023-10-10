package category

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var GetCategoriesBudgetValidator = validator.MustForm(map[string]validator.Validator{
	"app_meta": entity.AppMetaValidator(),
	"category_ids": &validator.Slice{
		Optional:  false,
		Validator: &validator.String{},
	},
	"budget_date": &validator.String{
		Optional:   false,
		Validators: []validator.StringFunc{entity.CheckDateStr},
	},
})

func (h *categoryHandler) GetCategoriesBudget(
	ctx context.Context,
	req *presenter.GetCategoriesBudgetRequest,
	res *presenter.GetCategoriesBudgetResponse,
) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.categoryUseCase.GetCategoriesBudget(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get categories budget, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
