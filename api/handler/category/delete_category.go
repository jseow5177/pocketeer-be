package category

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var DeleteCategoryValidator = validator.MustForm(map[string]validator.Validator{
	"category_id": &validator.String{
		Optional: false,
	},
})

func (h *categoryHandler) DeleteCategory(ctx context.Context, req *presenter.DeleteCategoryRequest, res *presenter.DeleteCategoryResponse) error {
	userID := util.GetUserIDFromCtx(ctx)

	useCaseRes, err := h.categoryUseCase.DeleteCategory(ctx, req.ToUseCaseReq(userID))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to delete category, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
