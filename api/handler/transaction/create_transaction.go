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

var CreateTransactionValidator = validator.MustForm(map[string]validator.Validator{
	"cat_id": &validator.String{
		Optional: false,
	},
	"amount": &validator.String{
		Optional:   false,
		Validators: []validator.StringFunc{entity.CheckAmount},
	},
	"transaction_type": &validator.Uint32{
		Optional:   false,
		Validators: []validator.Uint32Func{entity.CheckTransactionType},
	},
	"transaction_time": &validator.Uint64{
		Optional: false,
	},
})

func (h *TransactionHandler) CreateTransaction(ctx context.Context, req *presenter.CreateTransactionRequest, res *presenter.CreateTransactionResponse) error {
	var (
		userID = middleware.GetUserIDFromCtx(ctx)
		uc     = transaction.NewTransactionUseCase(h.categoryUseCase, h.transactionRepo)
	)

	t, err := uc.CreateTransaction(ctx, userID, req)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to create transaction, err: %v", err)
		return err
	}

	res.SetTransaction(t)

	return nil
}
