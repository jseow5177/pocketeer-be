package category

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/middleware"
	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var GetCategoryValidator = validator.MustForm(map[string]validator.Validator{
	"category_id": &validator.String{
		Optional: false,
	},
})

func (h *CategoryHandler) GetCategory(ctx context.Context, req *presenter.GetCategoryRequest, res *presenter.GetCategoryResponse) error {
	userID := middleware.GetUserIDFromCtx(ctx)

	c, err := h.categoryUseCase.GetCategory(ctx, userID, req)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get category, err: %v", err)
		return err
	}

	res.SetCategory(c)

	return nil
}
