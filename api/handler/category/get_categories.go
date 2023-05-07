package category

import (
	"context"

	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/data/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/usecase/category"
	"github.com/rs/zerolog/log"
)

var GetCategoriesValidator = validator.MustForm(map[string]validator.Validator{
	"cat_type": &validator.Uint32{
		Optional:   true,
		UnsetZero:  true,
		Validators: []validator.Uint32Func{entity.CheckCatType},
	},
})

func (h *CategoryHandler) GetCategories(ctx context.Context, req *presenter.GetCategoriesRequest, res *presenter.GetCategoriesResponse) error {
	var (
		uc = category.NewCategoryUseCase(h.categoryRepo)
		cf = req.ToCategoryFilter()
	)

	ecs, err := uc.GetCategories(ctx, cf)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get categories, err: %v", err)
		return err
	}

	res.SetCategories(ecs)

	return nil
}
