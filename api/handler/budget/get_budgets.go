package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
)

var GetBudgetsValidator = validator.MustForm(map[string]validator.Validator{})

func (h *budgetHandler) GetBudgets(ctx context.Context, req *presenter.GetBudgetsRequest, res *presenter.GetBudgetsResponse) error {
	return nil
}
