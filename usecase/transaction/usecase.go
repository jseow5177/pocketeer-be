package transaction

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/account"
	"github.com/jseow5177/pockteer-be/usecase/category"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var (
	ErrMismatchTransactionType = errors.New("mismatch transaction type")
	ErrInvalidCategoryIDs      = errors.New("invalid category_ids")
)

type transactionUseCase struct {
	categoryUseCase category.UseCase
	accountUseCase  account.UseCase
	transactionRepo repo.TransactionRepo
}

func NewTransactionUseCase(
	categoryUseCase category.UseCase,
	accountUseCase account.UseCase,
	transactionRepo repo.TransactionRepo,
) UseCase {
	return &transactionUseCase{
		categoryUseCase,
		accountUseCase,
		transactionRepo,
	}
}

func (uc *transactionUseCase) GetTransaction(ctx context.Context, req *GetTransactionRequest) (*GetTransactionResponse, error) {
	t, err := uc.transactionRepo.Get(ctx, req.ToTransactionFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transaction from repo, err: %v", err)
		return nil, err
	}

	cRes, err := uc.categoryUseCase.GetCategory(ctx, req.ToGetCategoryRequest(t.GetCategoryID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get category from repo, err: %v", err)
		return nil, err
	}

	acRes, err := uc.accountUseCase.GetAccount(ctx, req.ToGetAccountRequest(t.GetAccountID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get account from repo, err: %v", err)
		return nil, err
	}

	return &GetTransactionResponse{
		Transaction: &Transaction{
			Transaction: t,
			Category:    cRes.Category,
			Account:     acRes.Account,
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

	ucts := make([]*Transaction, len(ts))
	if c != nil {
		for i, t := range ts {
			ucts[i] = &Transaction{
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
			ucts[i] = &Transaction{
				Transaction: t,
				Category:    c.Category,
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}

	return &GetTransactionsResponse{
		Transactions: ucts,
		Paging:       req.Paging,
	}, nil
}

func (uc *transactionUseCase) CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*CreateTransactionResponse, error) {
	t := req.ToTransactionEntity()

	cRes, err := uc.categoryUseCase.GetCategory(ctx, req.ToGetCategoryRequest())
	if err != nil {
		return nil, err
	}

	if cRes.Category.GetCategoryType() != req.GetTransactionType() {
		return nil, ErrMismatchTransactionType
	}

	id, err := uc.transactionRepo.Create(ctx, t)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save new transaction to repo, err: %v", err)
		return nil, err
	}

	t.TransactionID = goutil.String(id)

	return &CreateTransactionResponse{
		Transaction: &Transaction{
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
	uct := tRes.Transaction

	nt := uc.getTransactionUpdates(uct.Transaction, req.ToTransactionEntity())
	if nt == nil {
		// no updates
		log.Ctx(ctx).Info().Msg("transaction has no updates")
		return &UpdateTransactionResponse{
			uct,
		}, nil
	}

	// check if category exists
	if nt.CategoryID != nil {
		cRes, err := uc.categoryUseCase.GetCategory(ctx, req.ToGetCategoryRequest())
		if err != nil {
			return nil, err
		}
		uct.Category = cRes.Category
	}

	tf := req.ToTransactionFilter()
	if err = uc.transactionRepo.Update(ctx, tf, nt); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save transaction updates to repo, err: %v", err)
		return nil, err
	}

	// merge
	goutil.MergeWithPtrFields(uct.Transaction, nt)

	return &UpdateTransactionResponse{
		uct,
	}, nil
}

func (uc *transactionUseCase) AggrTransactions(ctx context.Context, req *AggrTransactionsRequest) (*AggrTransactionsResponse, error) {
	tf := req.ToTransactionFilter(req.GetUserID())

	// default sum by category_id
	sumBy := "category_id"

	switch {
	case len(req.CategoryIDs) > 0:
		getCategoriesRes, err := uc.categoryUseCase.GetCategories(ctx, req.ToGetCategoriesRequest())
		if err != nil {
			return nil, err
		}

		if len(req.CategoryIDs) != len(getCategoriesRes.Categories) {
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

	if changes.Note != nil && changes.GetNote() != old.GetNote() {
		hasUpdates = true
		nt.Note = changes.Note
	}

	if !hasUpdates {
		return nil
	}

	return nt
}
