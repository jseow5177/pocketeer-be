package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var CreateBudgetsValidator = validator.MustForm(map[string]validator.Validator{
	"budgets": &validator.Slice{
		Optional:  false,
		MaxLen:    20,
		Validator: CreateBudgetValidator,
	},
})

func (h *budgetHandler) CreateBudgets(ctx context.Context, req *presenter.CreateBudgetsRequest, res *presenter.CreateBudgetsResponse) error {
	userID := util.GetUserIDFromCtx(ctx)

	useCaseRes, err := h.budgetUseCase.CreateBudgets(ctx, req.ToUseCaseReq(userID))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to create budgets, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
