package transaction

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var GetTransactionsValidator = validator.MustForm(map[string]validator.Validator{
	"category_id": &validator.String{
		Optional: true,
	},
	"account_id": &validator.String{
		Optional:  true,
		UnsetZero: true,
	},
	"category_ids": &validator.Slice{
		Optional:  true,
		Validator: new(validator.String),
	},
	"transaction_type": &validator.UInt32{
		Optional:   true,
		UnsetZero:  true,
		Validators: []validator.UInt32Func{entity.CheckTransactionType},
	},
	"transaction_time": validator.MustForm(map[string]validator.Validator{
		"gte": &validator.UInt64{
			Optional: false,
		},
		"lte": &validator.UInt64{
			Optional: false,
		},
	}),
	"paging": entity.PagingValidator(true),
})

func (h *transactionHandler) GetTransactions(ctx context.Context, req *presenter.GetTransactionsRequest, res *presenter.GetTransactionsResponse) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.transactionUseCase.GetTransactions(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transaction, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
