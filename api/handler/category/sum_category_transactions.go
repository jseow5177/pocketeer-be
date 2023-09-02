package category

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var SumCategoryTransactionsValidator = validator.MustForm(map[string]validator.Validator{
	"transaction_time": entity.RangeFilterValidator(true),
})

func (h *categoryHandler) SumCategoryTransactions(
	ctx context.Context,
	req *presenter.SumCategoryTransactionsRequest,
	res *presenter.SumCategoryTransactionsResponse,
) error {
	userID := util.GetUserIDFromCtx(ctx)

	useCaseRes, err := h.categoryUseCase.SumCategoryTransactions(ctx, req.ToUseCaseReq(userID))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to sum category transactions, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
