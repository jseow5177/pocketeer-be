package category

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
)

var UpdateCategoryValidator = validator.MustForm(map[string]validator.Validator{
	"category_id": &validator.String{
		Optional: false,
	},
	"category_name": &validator.String{
		Optional:  true,
		UnsetZero: true,
	},
})

func (h *categoryHandler) UpdateCategory(ctx context.Context, req *presenter.UpdateCategoryRequest, res *presenter.UpdateCategoryResponse) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.categoryUseCase.UpdateCategory(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to update category, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
