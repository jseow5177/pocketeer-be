package category

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var GetCategoryBudgetValidator = validator.MustForm(map[string]validator.Validator{
	"app_meta": entity.AppMetaValidator(),
	"category_id": &validator.String{
		Optional: false,
	},
	"budget_date": &validator.String{
		Optional:   false,
		Validators: []validator.StringFunc{entity.CheckDateStr},
	},
})

func (h *categoryHandler) GetCategoryBudget(
	ctx context.Context,
	req *presenter.GetCategoryBudgetRequest,
	res *presenter.GetCategoryBudgetResponse,
) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.categoryUseCase.GetCategoryBudget(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get category budget, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
