package transaction

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
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
	budgetRepo      repo.BudgetRepo
}

func NewTransactionUseCase(
	txMgr repo.TxMgr,
	categoryRepo repo.CategoryRepo,
	accountRepo repo.AccountRepo,
	transactionRepo repo.TransactionRepo,
	budgetRepo repo.BudgetRepo,
) UseCase {
	return &transactionUseCase{
		txMgr,
		categoryRepo,
		accountRepo,
		transactionRepo,
		budgetRepo,
	}
}

func (uc *transactionUseCase) GetTransaction(ctx context.Context, req *GetTransactionRequest) (*GetTransactionResponse, error) {
	t, err := uc.transactionRepo.Get(ctx, req.ToTransactionFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transaction from repo, err: %v", err)
		return nil, err
	}

	return &GetTransactionResponse{
		Transaction: t,
	}, nil
}

func (uc *transactionUseCase) GetTransactions(ctx context.Context, req *GetTransactionsRequest) (*GetTransactionsResponse, error) {
	ts, err := uc.transactionRepo.GetMany(ctx, req.ToTransactionFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transactions from repo, err: %v", err)
		return nil, err
	}

	return &GetTransactionsResponse{
		Transactions: ts,
		Paging:       req.Paging,
	}, nil
}

func (uc *transactionUseCase) DeleteTransaction(ctx context.Context, req *DeleteTransactionRequest) (*DeleteTransactionResponse, error) {
	t, err := uc.transactionRepo.Get(ctx, req.ToTransactionFilter())
	if err != nil && err != repo.ErrTransactionNotFound {
		log.Ctx(ctx).Error().Msgf("fail to get transaction from repo, err: %v", err)
		return nil, err
	}

	if err == repo.ErrTransactionNotFound {
		return new(DeleteTransactionResponse), nil
	}

	if err = uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		// mark transaction as deleted
		if err = uc.transactionRepo.Update(txCtx, req.ToTransactionFilter(), req.ToTransactionUpdate()); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to mark transaction as deleted, err: %v", err)
			return err
		}

		if t.GetAmount() != 0 {
			// get account
			ac, err := uc.accountRepo.Get(ctx, req.ToAccountFilter(t))
			if err != nil {
				log.Ctx(txCtx).Error().Msgf("fail to get account from repo, err: %v", err)
				return err
			}

			// reset account balance
			newBalance := ac.GetBalance() - t.GetAmount()
			nac, hasUpdate, err := ac.Update(entity.NewAccountUpdate(
				entity.WithUpdateAccountBalance(goutil.Float64(newBalance)),
			))
			if err != nil {
				return err
			}

			if hasUpdate {
				if err := uc.accountRepo.Update(txCtx, req.ToAccountFilter(t), nac); err != nil {
					log.Ctx(txCtx).Error().Msgf("fail to update account balance, err: %v", err)
					return err
				}
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return new(DeleteTransactionResponse), nil
}

func (uc *transactionUseCase) CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*CreateTransactionResponse, error) {
	t := req.ToTransactionEntity()

	c, err := uc.categoryRepo.Get(ctx, req.ToCategoryFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get category from repo, err: %v", err)
		return nil, err
	}

	_, err = t.CanTransactionUnderCategory(c)
	if err != nil {
		return nil, err
	}

	ac, err := uc.accountRepo.Get(ctx, req.ToAccountFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get account from repo, err: %v", err)
		return nil, err
	}

	_, err = t.CanTransactionUnderAccount(ac)
	if err != nil {
		return nil, err
	}

	if err := uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		// create transaction
		_, err := uc.transactionRepo.Create(txCtx, t)
		if err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to save new transaction to repo, err: %v", err)
			return err
		}

		// update account balance
		newBalance := ac.GetBalance() + t.GetAmount()
		nac, _, err := ac.Update(entity.NewAccountUpdate(
			entity.WithUpdateAccountBalance(goutil.Float64(newBalance)),
		))
		if err != nil {
			return err
		}
		if err := uc.accountRepo.Update(txCtx, req.ToAccountFilter(), nac); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to update account balance, err: %v", err)
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &CreateTransactionResponse{
		Transaction: t,
	}, nil
}

func (uc *transactionUseCase) UpdateTransaction(ctx context.Context, req *UpdateTransactionRequest) (*UpdateTransactionResponse, error) {
	t, err := uc.transactionRepo.Get(ctx, req.ToTransactionFilter())
	if err != nil {
		return nil, err
	}
	oldAmount := t.GetAmount()

	tu, hasUpdate := t.Update(req.ToTransactionUpdate())
	if !hasUpdate {
		log.Ctx(ctx).Info().Msg("transaction has no updates")
		return &UpdateTransactionResponse{
			Transaction: t,
		}, nil
	}

	ac, err := uc.accountRepo.Get(ctx, req.ToAccountFilter(t.GetAccountID()))
	if err != nil {
		log.Ctx(ctx).Info().Msgf("fail to get account from repo, err: %v", err)
		return nil, err
	}

	if err = uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		// save updates
		if err = uc.transactionRepo.Update(txCtx, req.ToTransactionFilter(), tu); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to save transaction updates to repo, err: %v", err)
			return err
		}

		// update balance
		if tu.Amount != nil {
			newBalance := ac.GetBalance() + (tu.GetAmount() - oldAmount)
			nac, hasUpdate, err := ac.Update(entity.NewAccountUpdate(
				entity.WithUpdateAccountBalance(goutil.Float64(newBalance)),
			))
			if err != nil {
				return err
			}

			if hasUpdate {
				if err := uc.accountRepo.Update(txCtx, req.ToAccountFilter(t.GetAccountID()), nac); err != nil {
					log.Ctx(txCtx).Error().Msgf("fail to update account balance, err: %v", err)
					return err
				}
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &UpdateTransactionResponse{
		Transaction: t,
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

	aggrs, err := uc.transactionRepo.CalcTotalAmount(ctx, sumBy, tf)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to sum transactions by %s, err: %v", sumBy, err)
		return nil, err
	}

	results := make(map[string]*Aggr)
	for _, aggr := range aggrs {
		results[aggr.GetGroupBy()] = &Aggr{
			Sum: aggr.TotalAmount,
		}
	}

	return &AggrTransactionsResponse{
		Results: results,
	}, nil
}
