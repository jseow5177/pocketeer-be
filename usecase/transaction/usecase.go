package transaction

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidCategoryIDs = errors.New("invalid category_ids")
)

type transactionUseCase struct {
	txMgr           repo.TxMgr
	categoryRepo    repo.CategoryRepo
	accountRepo     repo.AccountRepo
	transactionRepo repo.TransactionRepo
}

func NewTransactionUseCase(
	txMgr repo.TxMgr,
	categoryRepo repo.CategoryRepo,
	accountRepo repo.AccountRepo,
	transactionRepo repo.TransactionRepo,
) UseCase {
	return &transactionUseCase{
		txMgr,
		categoryRepo,
		accountRepo,
		transactionRepo,
	}
}

func (uc *transactionUseCase) GetTransaction(ctx context.Context, req *GetTransactionRequest) (*GetTransactionResponse, error) {
	t, err := uc.transactionRepo.Get(ctx, req.ToTransactionFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transaction from repo, err: %v", err)
		return nil, err
	}

	return &GetTransactionResponse{
		t,
	}, nil
}

func (uc *transactionUseCase) GetTransactions(ctx context.Context, req *GetTransactionsRequest) (*GetTransactionsResponse, error) {
	ts, err := uc.transactionRepo.GetMany(ctx, req.ToTransactionFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transactions from repo, err: %v", err)
		return nil, err
	}

	return &GetTransactionsResponse{
		ts,
		req.Paging,
	}, nil
}

func (uc *transactionUseCase) CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*CreateTransactionResponse, error) {
	t := req.ToTransactionEntity()

	c, err := uc.categoryRepo.Get(ctx, req.ToCategoryFilter())
	if err != nil {
		return nil, err
	}

	_, err = t.CanTransactionUnderCategory(c)
	if err != nil {
		return nil, err
	}

	_, err = uc.accountRepo.Get(ctx, req.ToAccountFilter())
	if err != nil {
		return nil, err
	}

	if err := uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		// create transaction
		_, err := uc.transactionRepo.Create(ctx, t)
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to save new transaction to repo, err: %v", err)
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &CreateTransactionResponse{
		t,
	}, nil
}

func (uc *transactionUseCase) UpdateTransaction(ctx context.Context, req *UpdateTransactionRequest) (*UpdateTransactionResponse, error) {
	t, err := uc.transactionRepo.Get(ctx, req.ToTransactionFilter())
	if err != nil {
		return nil, err
	}

	nt := t.GetUpdates(req.ToTransactionUpdate(), true)
	if nt == nil {
		log.Ctx(ctx).Info().Msg("transaction has no updates")
		return &UpdateTransactionResponse{
			t,
		}, nil
	}

	// check if category exists
	if req.CategoryID != nil {
		_, err := uc.categoryRepo.Get(ctx, req.ToCategoryFilter())
		if err != nil {
			log.Ctx(ctx).Info().Msgf("fail to get category from repo, categoryID: %v, err: %v", req.GetCategoryID(), err)
			return nil, err
		}
	}

	// check if account exists
	if req.AccountID != nil {
		_, err := uc.accountRepo.Get(ctx, req.ToAccountFilter(t.GetAccountID()))
		if err != nil {
			log.Ctx(ctx).Info().Msgf("fail to get account from repo, accountID: %v, err: %v", t.GetAccountID(), err)
			return nil, err
		}
	}

	if err = uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		// save updates
		if err = uc.transactionRepo.Update(ctx, req.ToTransactionFilter(), nt); err != nil {
			log.Ctx(ctx).Error().Msgf("fail to save transaction updates to repo, err: %v", err)
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &UpdateTransactionResponse{
		t,
	}, nil
}

func (uc *transactionUseCase) AggrTransactions(ctx context.Context, req *AggrTransactionsRequest) (*AggrTransactionsResponse, error) {
	tf := req.ToTransactionFilter(req.GetUserID())

	// default sum by category_id
	sumBy := "category_id"

	switch {
	case len(req.CategoryIDs) > 0:
		categories, err := uc.categoryRepo.GetMany(ctx, req.ToCategoryFilter())
		if err != nil {
			return nil, err
		}

		if len(req.CategoryIDs) != len(categories) {
			return nil, ErrInvalidCategoryIDs
		}
	case len(req.TransactionTypes) > 0:
		sumBy = "transaction_type"
	}

	res, err := uc.transactionRepo.Sum(ctx, sumBy, tf)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to sum transactions by %s, err: %v", sumBy, err)
		return nil, err
	}

	results := make(map[string]*Aggr)
	for v, s := range res {
		results[v] = &Aggr{
			Sum: goutil.Float64(s),
		}
	}

	return &AggrTransactionsResponse{
		Results: results,
	}, nil
}
