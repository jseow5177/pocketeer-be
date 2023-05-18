package transaction

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/middleware"
	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/usecase/transaction"
	"github.com/rs/zerolog/log"
)

var GetTransactionsValidator = validator.MustForm(map[string]validator.Validator{
	"category_id": &validator.String{
		Optional:  true,
		UnsetZero: true,
	},
	"transaction_type": &validator.UInt32{
		Optional:   true,
		UnsetZero:  true,
		Validators: []validator.UInt32Func{entity.CheckTransactionType},
	},
	"transaction_time": entity.UInt64FilterValidator(true),
	"paging":           entity.PagingValidator(true),
})

func (h *TransactionHandler) GetTransactions(ctx context.Context, req *presenter.GetTransactionsRequest, res *presenter.GetTransactionsResponse) error {
	var (
		userID = middleware.GetUserIDFromCtx(ctx)
		uc     = transaction.NewTransactionUseCase(h.categoryUseCase, h.transactionRepo)
	)

	if req.Paging == nil {
		req.Paging = new(presenter.Paging)
	}

	if req.Paging.Limit == nil {
		req.Paging.Limit = goutil.Uint32(config.DefaultPagingLimit)
	}
	if req.Paging.Page == nil {
		req.Paging.Page = goutil.Uint32(config.MinPagingPage)
	}

	ts, err := uc.GetTransactions(ctx, userID, req)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transaction, err: %v", err)
		return err
	}

	res.SetTransactions(ts)
	res.SetPaging(req.Paging)

	return nil
}
