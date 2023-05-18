package transaction

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/middleware"
	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/usecase/transaction"
	"github.com/rs/zerolog/log"
)

var GetTransactionValidator = validator.MustForm(map[string]validator.Validator{
	"transaction_id": &validator.String{
		Optional: false,
	},
})

func (h *transactionHandler) GetTransaction(ctx context.Context, req *presenter.GetTransactionRequest, res *presenter.GetTransactionResponse) error {
	var (
		userID = middleware.GetUserIDFromCtx(ctx)
		uc     = transaction.NewTransactionUseCase(h.categoryUseCase, h.transactionRepo)
	)

	t, err := uc.GetTransaction(ctx, userID, req)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transaction, err: %v", err)
		return err
	}

	res.SetTransaction(t)

	return nil
}
