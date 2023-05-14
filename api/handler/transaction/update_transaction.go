package transaction

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/middleware"
	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/data/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/usecase/transaction"
	"github.com/rs/zerolog/log"
)

var UpdateTransactionValidator = validator.MustForm(map[string]validator.Validator{
	"transaction_id": &validator.String{
		Optional: false,
	},
	"category_id": &validator.String{
		Optional:  true,
		UnsetZero: true,
	},
	"amount": &validator.String{
		Optional:   true,
		UnsetZero:  true,
		Validators: []validator.StringFunc{entity.CheckAmount},
	},
	"transaction_type": &validator.UInt32{
		Optional:   true,
		UnsetZero:  true,
		Validators: []validator.UInt32Func{entity.CheckTransactionType},
	},
	"transaction_time": &validator.UInt64{
		Optional:  true,
		UnsetZero: true,
	},
})

func (h *TransactionHandler) UpdateTransaction(ctx context.Context, req *presenter.UpdateTransactionRequest, res *presenter.UpdateTransactionResponse) error {
	var (
		userID = middleware.GetUserIDFromCtx(ctx)
		uc     = transaction.NewTransactionUseCase(h.categoryUseCase, h.transactionRepo)
	)

	t, err := uc.UpdateTransaction(ctx, userID, req)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to update transaction, err: %v", err)
		return err
	}

	res.SetTransaction(t)

	return nil
}
