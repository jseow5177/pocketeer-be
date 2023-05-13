package transaction

import (
	"context"
	"time"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/data/presenter"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/pkg/errutil"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/category"
	"github.com/rs/zerolog/log"
)

type TransactionUseCase struct {
	categoryRepo    category.UseCase
	transactionRepo repo.TransactionRepo
}

func NewTransactionUseCase(categoryUseCase category.UseCase, transactionRepo repo.TransactionRepo) UseCase {
	return &TransactionUseCase{
		categoryRepo:    categoryUseCase,
		transactionRepo: transactionRepo,
	}
}

func (uc *TransactionUseCase) GetTransaction(ctx context.Context, userID string, req *presenter.GetTransactionRequest) (*entity.Transaction, error) {
	tf := req.ToTransactionFilter(userID)

	t, err := uc.transactionRepo.Get(ctx, tf)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transaction from repo, err: %v", err)
		if err == errutil.ErrNotFound {
			return nil, errutil.NotFoundError(err)
		}
		return nil, err
	}

	return t, nil
}

func (uc *TransactionUseCase) GetTransactions(ctx context.Context, userID string, req *presenter.GetTransactionsRequest) ([]*entity.Transaction, error) {
	return nil, nil
}

func (uc *TransactionUseCase) CreateTransaction(ctx context.Context, userID string, req *presenter.CreateTransactionRequest) (*entity.Transaction, error) {
	var (
		t   = req.ToTransactionEntity(userID)
		now = uint64(time.Now().Unix())
	)

	if _, err := uc.categoryRepo.GetCategory(ctx, userID, req.ToGetCategoryRequest()); err != nil {
		return nil, err
	}

	t.CreateTime = goutil.Uint64(now)
	t.UpdateTime = goutil.Uint64(now)

	// Standardize amount to two decimal places
	if err := t.StandardizeAmount(config.AmountDecimalPlaces); err != nil {
		log.Ctx(ctx).Error().Msgf("invalid amount, amount: %v, err: %v", t.GetAmount(), err)
		return nil, errutil.BadRequestError(err)
	}

	id, err := uc.transactionRepo.Create(ctx, t)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save new transaction to repo, err: %v", err)
		return nil, err
	}

	t.TransactionID = goutil.String(id)

	return t, nil
}
