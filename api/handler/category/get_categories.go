package category

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var GetCategoriesValidator = validator.MustForm(map[string]validator.Validator{
	"category_type": &validator.UInt32{
		Optional:   true,
		UnsetZero:  true,
		Validators: []validator.UInt32Func{entity.CheckCategoryType},
	},
	"category_ids": &validator.Slice{
		Optional:  true,
		Validator: &validator.String{},
	},
})

func (h *categoryHandler) GetCategories(ctx context.Context, req *presenter.GetCategoriesRequest, res *presenter.GetCategoriesResponse) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.categoryUseCase.GetCategories(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get categories, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
