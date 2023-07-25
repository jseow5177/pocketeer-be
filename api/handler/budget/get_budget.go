package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
)

var GetBudgetValidator = validator.MustForm(map[string]validator.Validator{})

func (h *budgetHandler) GetBudget(ctx context.Context, req *presenter.GetBudgetRequest, res *presenter.GetBudgetResponse) error {
	return nil
}
