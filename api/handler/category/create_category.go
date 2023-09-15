package category

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
)

var CreateCategoryValidator = validator.MustForm(map[string]validator.Validator{
	"category_name": &validator.String{
		Optional: false,
	},
	"category_type": &validator.UInt32{
		Optional:   false,
		Validators: []validator.UInt32Func{entity.CheckCategoryType},
	},
})

func (h *categoryHandler) CreateCategory(ctx context.Context, req *presenter.CreateCategoryRequest, res *presenter.CreateCategoryResponse) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.categoryUseCase.CreateCategory(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to create category, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
