package category

import (
	"context"

	"github.com/jseow5177/pockteer-be/model"
	"github.com/jseow5177/pockteer-be/pkg/validator"
)

var CreateCategoryValidator = validator.MustForm(map[string]validator.Validator{})

func CreateCategory(ctx context.Context, req *model.CreateCategoryRequest, res *model.CreateCategoryResponse) error {
	return nil
}
