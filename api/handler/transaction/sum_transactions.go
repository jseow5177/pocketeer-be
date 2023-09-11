package transaction

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var SumTransactionsValidator = validator.MustForm(map[string]validator.Validator{
	"transaction_time": validator.MustForm(map[string]validator.Validator{
		"gte": &validator.UInt64{
			Optional: false,
		},
		"lte": &validator.UInt64{
			Optional: false,
		},
	}),
	"transaction_type": &validator.UInt32{
		Optional:   true,
		Validators: []validator.UInt32Func{entity.CheckTransactionType},
	},
})

func (h *transactionHandler) SumTransactions(
	ctx context.Context,
	req *presenter.SumTransactionsRequest,
	res *presenter.SumTransactionsResponse,
) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.transactionUseCase.SumTransactions(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to sum transaction, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
