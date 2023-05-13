package transaction

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/middleware"
	"github.com/jseow5177/pockteer-be/data/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/usecase/transaction"
	"github.com/rs/zerolog/log"
)

var GetTransactionsValidator = validator.MustForm(map[string]validator.Validator{})

func (h *TransactionHandler) GetTransactions(ctx context.Context, req *presenter.GetTransactionsRequest, res *presenter.GetTransactionsResponse) error {
	var (
		userID = middleware.GetUserIDFromCtx(ctx)
		uc     = transaction.NewTransactionUseCase(h.categoryUseCase, h.transactionRepo)
	)

	ts, err := uc.GetTransactions(ctx, userID, req)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transaction, err: %v", err)
		return err
	}

	res.SetTransactions(ts)

	return nil
}
