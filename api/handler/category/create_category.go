package category

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/data/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/usecase/category"
)

var CreateCategoryValidator = validator.MustForm(map[string]validator.Validator{
	"cat_name": &validator.String{
		Optional: false,
	},
	"cat_type": &validator.Uint32{
		Optional:   false,
		Validators: []validator.Uint32Func{entity.CheckCatType},
	},
})

func (h *CategoryHandler) CreateCategory(ctx context.Context, req *presenter.CreateCategoryRequest, res *presenter.CreateCategoryResponse) error {
	var (
		uc = category.NewCategoryUseCase(h.categoryRepo)
		c  = req.ToCategoryEntity()
	)

	if err := uc.CreateCategory(ctx, c); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to create category, err: %v", err)
		return err
	}

	res.ToCategoryPresenter(c)

	return nil
}
