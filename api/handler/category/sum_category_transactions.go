package category

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var SumCategoryTransactionsValidator = validator.MustForm(map[string]validator.Validator{
	"transaction_time": entity.RangeFilterValidator(true),
	"transaction_type": &validator.UInt32{
		Optional:   false,
		Validators: []validator.UInt32Func{entity.CheckCategoryType},
	},
})

func (h *categoryHandler) SumCategoryTransactions(
	ctx context.Context,
	req *presenter.SumCategoryTransactionsRequest,
	res *presenter.SumCategoryTransactionsResponse,
) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.categoryUseCase.SumCategoryTransactions(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to sum category transactions, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
