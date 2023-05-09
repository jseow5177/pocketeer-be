package category

import (
	"context"

	"github.com/jseow5177/pockteer-be/data/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/usecase/category"
	"github.com/rs/zerolog/log"
)

var GetCategoryValidator = validator.MustForm(map[string]validator.Validator{
	"cat_id": &validator.String{
		Optional: false,
	},
})

func (h *CategoryHandler) GetCategory(ctx context.Context, req *presenter.GetCategoryRequest, res *presenter.GetCategoryResponse) error {
	var (
		uc = category.NewCategoryUseCase(h.categoryRepo)
		cf = req.ToCategoryFilter()
	)

	ec, err := uc.GetCategory(ctx, cf)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get category, err: %v", err)
		return err
	}

	res.SetCategory(ec)

	return nil
}
