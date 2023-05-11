package category

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/middleware"
	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/data/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var GetCategoriesValidator = validator.MustForm(map[string]validator.Validator{
	"cat_type": &validator.Uint32{
		Optional:   true,
		UnsetZero:  true,
		Validators: []validator.Uint32Func{entity.CheckCategoryType},
	},
})

func (h *CategoryHandler) GetCategories(ctx context.Context, req *presenter.GetCategoriesRequest, res *presenter.GetCategoriesResponse) error {
	userID := middleware.GetUserIDFromCtx(ctx)

	cs, err := h.categoryUseCase.GetCategories(ctx, userID, req)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get categories, err: %v", err)
		return err
	}

	res.SetCategories(cs)

	return nil
}
