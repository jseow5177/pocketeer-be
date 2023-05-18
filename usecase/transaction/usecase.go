package transaction

import (
	"context"
	"errors"
	"time"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/errutil"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/category"
	"github.com/rs/zerolog/log"
)

var (
	ErrTransactionNotFound     = errors.New("transaction not found")
	ErrMismatchTransactionType = errors.New("mismatch transaction type")
)

type TransactionUseCase struct {
	categoryUseCase category.UseCase
	transactionRepo repo.TransactionRepo
}

func NewTransactionUseCase(categoryUseCase category.UseCase, transactionRepo repo.TransactionRepo) UseCase {
	return &TransactionUseCase{
		categoryUseCase: categoryUseCase,
		transactionRepo: transactionRepo,
	}
}

func (uc *TransactionUseCase) GetTransaction(ctx context.Context, userID string, req *presenter.GetTransactionRequest) (*entity.Transaction, error) {
	tf := req.ToTransactionFilter(userID)

	t, err := uc.transactionRepo.Get(ctx, tf)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transaction from repo, err: %v", err)
		if err == errutil.ErrNotFound {
			return nil, errutil.NotFoundError(ErrTransactionNotFound)
		}
		return nil, err
	}

	return t, nil
}

func (uc *TransactionUseCase) GetTransactions(ctx context.Context, userID string, req *presenter.GetTransactionsRequest) ([]*entity.Transaction, error) {
	tf := req.ToTransactionFilter(userID)

	if _, err := uc.categoryUseCase.GetCategory(ctx, userID, req.ToGetCategoryRequest()); err != nil {
		return nil, err
	}

	ts, err := uc.transactionRepo.GetMany(ctx, tf)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transactions from repo, err: %v", err)
		return nil, err
	}

	return ts, nil
}

func (uc *TransactionUseCase) CreateTransaction(ctx context.Context, userID string, req *presenter.CreateTransactionRequest) (*entity.Transaction, error) {
	var (
		t   = req.ToTransactionEntity(userID)
		now = uint64(time.Now().Unix())
	)

	c, err := uc.categoryUseCase.GetCategory(ctx, userID, req.ToGetCategoryRequest())
	if err != nil {
		return nil, err
	}

	if c.GetCategoryType() != req.GetTransactionType() {
		return nil, errutil.ValidationError(ErrMismatchTransactionType)
	}

	t.CreateTime = goutil.Uint64(now)
	t.UpdateTime = goutil.Uint64(now)

	id, err := uc.transactionRepo.Create(ctx, t)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save new transaction to repo, err: %v", err)
		return nil, err
	}

	t.TransactionID = goutil.String(id)

	return t, nil
}

func (uc *TransactionUseCase) UpdateTransaction(ctx context.Context, userID string, req *presenter.UpdateTransactionRequest) (*entity.Transaction, error) {
	t, err := uc.GetTransaction(ctx, userID, req.ToGetTransactionRequest())
	if err != nil {
		return nil, err
	}

	nt := uc.getTransactionUpdates(t, req.ToTransactionEntity())
	if nt == nil {
		// no updates
		log.Ctx(ctx).Info().Msg("transaction has no updates")
		return t, nil
	}

	tf := req.ToTransactionFilter(userID)
	if err = uc.transactionRepo.Update(ctx, tf, nt); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save transaction updates to repo, err: %v", err)
		return nil, err
	}

	// merge
	goutil.MergeWithPtrFields(t, nt)

	return t, nil
}

func (uc *TransactionUseCase) getTransactionUpdates(old, changes *entity.Transaction) *entity.Transaction {
	var hasUpdates bool

	nt := &entity.Transaction{
		UpdateTime: goutil.Uint64(uint64(time.Now().Unix())),
	}

	if changes.CategoryID != nil && changes.GetCategoryID() != old.GetCategoryID() {
		hasUpdates = true
		nt.CategoryID = changes.CategoryID
	}

	if changes.Amount != nil && changes.GetAmount() != old.GetAmount() {
		hasUpdates = true
		nt.Amount = changes.Amount
	}

	if changes.TransactionType != nil && changes.GetTransactionType() != old.GetTransactionType() {
		hasUpdates = true
		nt.TransactionType = changes.TransactionType
	}

	if changes.TransactionTime != nil && changes.GetTransactionTime() != old.GetTransactionTime() {
		hasUpdates = true
		nt.TransactionTime = changes.TransactionTime
	}

	if !hasUpdates {
		return nil
	}

	return nt
}
