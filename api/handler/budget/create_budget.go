package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/model"
	"github.com/jseow5177/pockteer-be/pkg/validator"
)

var CreateBudgetValidator = validator.MustForm(map[string]validator.Validator{})

func CreateBudget(ctx context.Context, req *model.CreateBudgetRequest, res *model.CreateBudgetResponse) error {
	return nil
}
