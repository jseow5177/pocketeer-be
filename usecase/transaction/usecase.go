package transaction

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/errutil"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/category"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var (
	ErrMismatchTransactionType = errors.New("mismatch transaction type")
)

type transactionUseCase struct {
	categoryUseCase category.UseCase
	transactionRepo repo.TransactionRepo
}

func NewTransactionUseCase(categoryUseCase category.UseCase, transactionRepo repo.TransactionRepo) UseCase {
	return &transactionUseCase{
		categoryUseCase: categoryUseCase,
		transactionRepo: transactionRepo,
	}
}

func (uc *transactionUseCase) GetTransaction(ctx context.Context, req *GetTransactionRequest) (*GetTransactionResponse, error) {
	t, err := uc.transactionRepo.Get(ctx, req.ToTransactionFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transaction from repo, err: %v", err)
		if err == repo.ErrTransactionNotFound {
			return nil, errutil.NotFoundError(err)
		}
		return nil, err
	}

	cRes, err := uc.categoryUseCase.GetCategory(ctx, req.ToGetCategoryRequest(t.GetCategoryID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get category from repo, err: %v", err)
		return nil, err
	}

	return &GetTransactionResponse{
		TransactionWithCategory: &TransactionWithCategory{
			Transaction: t,
			Category:    cRes.Category,
		},
	}, nil
}

func (uc *transactionUseCase) GetTransactions(ctx context.Context, req *GetTransactionsRequest) (*GetTransactionsResponse, error) {
	var c *entity.Category
	if req.CategoryID != nil {
		cRes, err := uc.categoryUseCase.GetCategory(ctx, req.ToGetCategoryRequest())
		if err != nil {
			return nil, err
		}
		c = cRes.Category
	}

	tf := req.ToTransactionFilter()
	ts, err := uc.transactionRepo.GetMany(ctx, tf)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transactions from repo, err: %v", err)
		return nil, err
	}

	tswc := make([]*TransactionWithCategory, len(ts))
	if c != nil {
		for i, t := range ts {
			tswc[i] = &TransactionWithCategory{
				Transaction: t,
				Category:    c,
			}
		}
	} else {
		if err := util.ParallelizeWork(ctx, len(ts), 50, func(ctx context.Context, i int) error {
			t := ts[i]
			c, err := uc.categoryUseCase.GetCategory(ctx, &category.GetCategoryRequest{
				UserID:     req.UserID,
				CategoryID: t.CategoryID,
			})
			if err != nil {
				log.Ctx(ctx).Error().Msgf("fail to get category of transaction, transactionID: %v, err: %v",
					t.GetTransactionID(), err)
				return err
			}
			tswc[i] = &TransactionWithCategory{
				Transaction: t,
				Category:    c.Category,
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}

	return &GetTransactionsResponse{
		TransactionsWithCategory: tswc,
		Paging:                   req.Paging,
	}, nil
}

func (uc *transactionUseCase) CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*CreateTransactionResponse, error) {
	t := req.ToTransactionEntity()

	cRes, err := uc.categoryUseCase.GetCategory(ctx, req.ToGetCategoryRequest())
	if err != nil {
		return nil, err
	}

	if cRes.Category.GetCategoryType() != req.GetTransactionType() {
		return nil, errutil.ValidationError(ErrMismatchTransactionType)
	}

	id, err := uc.transactionRepo.Create(ctx, t)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save new transaction to repo, err: %v", err)
		return nil, err
	}

	t.TransactionID = goutil.String(id)

	return &CreateTransactionResponse{
		TransactionWithCategory: &TransactionWithCategory{
			Transaction: t,
			Category:    cRes.Category,
		},
	}, nil
}

func (uc *transactionUseCase) UpdateTransaction(ctx context.Context, req *UpdateTransactionRequest) (*UpdateTransactionResponse, error) {
	tRes, err := uc.GetTransaction(ctx, req.ToGetTransactionRequest())
	if err != nil {
		return nil, err
	}
	twc := tRes.TransactionWithCategory

	nt := uc.getTransactionUpdates(tRes.Transaction, req.ToTransactionEntity())
	if nt == nil {
		// no updates
		log.Ctx(ctx).Info().Msg("transaction has no updates")
		return &UpdateTransactionResponse{
			twc,
		}, nil
	}

	// check if category exists
	if nt.CategoryID != nil {
		cRes, err := uc.categoryUseCase.GetCategory(ctx, req.ToGetCategoryRequest())
		if err != nil {
			return nil, err
		}
		twc.Category = cRes.Category
	}

	tf := req.ToTransactionFilter()
	if err = uc.transactionRepo.Update(ctx, tf, nt); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save transaction updates to repo, err: %v", err)
		return nil, err
	}

	// merge
	goutil.MergeWithPtrFields(twc.Transaction, nt)

	return &UpdateTransactionResponse{
		twc,
	}, nil
}

func (uc *transactionUseCase) getTransactionUpdates(old, changes *entity.Transaction) *entity.Transaction {
	var hasUpdates bool

	nt := new(entity.Transaction)

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
