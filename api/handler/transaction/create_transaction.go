package transaction

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/middleware"
	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/usecase/transaction"
	"github.com/rs/zerolog/log"
)

var CreateTransactionValidator = validator.MustForm(map[string]validator.Validator{
	"category_id": &validator.String{
		Optional: false,
	},
	"amount": &validator.String{
		Optional:   false,
		Validators: []validator.StringFunc{entity.CheckAmount},
	},
	"transaction_type": &validator.UInt32{
		Optional:   false,
		Validators: []validator.UInt32Func{entity.CheckTransactionType},
	},
	"transaction_time": &validator.UInt64{
		Optional: false,
	},
})

func (h *transactionHandler) CreateTransaction(ctx context.Context, req *presenter.CreateTransactionRequest, res *presenter.CreateTransactionResponse) error {
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
