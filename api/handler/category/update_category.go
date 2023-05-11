package category

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/jseow5177/pockteer-be/data/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/usecase/category"
)

var UpdateCategoryValidator = validator.MustForm(map[string]validator.Validator{
	"cat_id": &validator.String{
		Optional: false,
	},
	"cat_name": &validator.String{
		Optional: true,
	},
})

func (h *CategoryHandler) UpdateCategory(ctx context.Context, req *presenter.UpdateCategoryRequest, res *presenter.UpdateCategoryResponse) error {
	uc := category.NewCategoryUseCase(h.categoryRepo)

	c, err := uc.UpdateCategory(ctx, req)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to update category, err: %v", err)
		return err
	}

	res.SetCategory(c)

	return nil
}
