package category

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/jseow5177/pockteer-be/api/middleware"
	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/data/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
)

var CreateCategoryValidator = validator.MustForm(map[string]validator.Validator{
	"cat_name": &validator.String{
		Optional: false,
	},
	"cat_type": &validator.Uint32{
		Optional:   false,
		Validators: []validator.Uint32Func{entity.CheckCategoryType},
	},
})

func (h *CategoryHandler) CreateCategory(ctx context.Context, req *presenter.CreateCategoryRequest, res *presenter.CreateCategoryResponse) error {
	userID := middleware.GetUserIDFromCtx(ctx)

	c, err := h.categoryUseCase.CreateCategory(ctx, userID, req)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to create category, err: %v", err)
		return err
	}

	res.SetCategory(c)

	return nil
}
