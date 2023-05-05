package category

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
)

var CreateCategoryValidator = validator.MustForm(map[string]validator.Validator{
	"cat_name": &validator.String{
		Optional: false,
	},
})

func (h *CategoryHandler) CreateCategory(ctx context.Context, req *presenter.CreateCategoryRequest, res *presenter.CreateCategoryResponse) error {
	return nil
}
