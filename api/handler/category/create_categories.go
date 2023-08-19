package category

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var CreateCategoriesValidator = validator.MustForm(map[string]validator.Validator{
	"categories": &validator.Slice{
		Optional:  false,
		MaxLen:    10,
		Validator: CreateCategoryValidator,
	},
})

func (h *categoryHandler) CreateCategories(ctx context.Context, req *presenter.CreateCategoriesRequest, res *presenter.CreateCategoriesResponse) error {
	userID := util.GetUserIDFromCtx(ctx)

	useCaseRes, err := h.categoryUseCase.CreateCategories(ctx, req.ToUseCaseReq(userID))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to create categories, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
