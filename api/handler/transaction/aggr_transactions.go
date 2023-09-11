package transaction

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var AggrTransactionsValidator = validator.MustForm(map[string]validator.Validator{
	"transaction_time": validator.MustForm(map[string]validator.Validator{
		"gte": &validator.UInt64{
			Optional: false,
		},
		"lte": &validator.UInt64{
			Optional: false,
		},
	}),
	"category_ids": &validator.Slice{
		Optional:  true,
		Validator: &validator.String{},
	},
	"budget_ids": &validator.Slice{
		Optional:  true,
		Validator: &validator.String{},
	},
	"transaction_types": &validator.Slice{
		Optional: true,
		Validator: &validator.UInt32{
			Validators: []validator.UInt32Func{entity.CheckTransactionType},
		},
	},
})

func (h *transactionHandler) AggrTransactions(
	ctx context.Context,
	req *presenter.AggrTransactionsRequest,
	res *presenter.AggrTransactionsResponse,
) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.transactionUseCase.AggrTransactions(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to aggr transaction, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
