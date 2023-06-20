package aggr

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/usecase/account"
	"github.com/jseow5177/pockteer-be/usecase/budget"
	"github.com/jseow5177/pockteer-be/usecase/category"
	"github.com/jseow5177/pockteer-be/usecase/transaction"
	"github.com/jseow5177/pockteer-be/util"
)

type aggrUseCase struct {
	budgetUseCase      budget.UseCase
	categoryUseCase    category.UseCase
	accountUseCase     account.UseCase
	transactionUseCase transaction.UseCase
}

func NewAggrUseCase(
	budgetUseCase budget.UseCase,
	categoryUseCase category.UseCase,
	accountUseCase account.UseCase,
	transactionUseCase transaction.UseCase,
) UseCase {
	return &aggrUseCase{
		budgetUseCase,
		categoryUseCase,
		accountUseCase,
		transactionUseCase,
	}
}

func (uc *aggrUseCase) GetTransaction(ctx context.Context, req *GetTransactionRequest) (*GetTransactionResponse, error) {
	tRes, err := uc.transactionUseCase.GetTransaction(ctx, req.ToGetTransactionRequest())
	if err != nil {
		return nil, err
	}

	t := tRes.Transaction

	cRes, err := uc.categoryUseCase.GetCategory(ctx, req.ToGetCategoryRequest(t.GetCategoryID()))
	if err != nil {
		return nil, err
	}

	acRes, err := uc.accountUseCase.GetAccount(ctx, req.ToGetAccountRequest(t.GetAccountID()))
	if err != nil {
		return nil, err
	}

	return &GetTransactionResponse{
		TransactionObj: &TransactionObj{
			t,
			cRes.Category,
			acRes.Account,
		},
	}, nil
}

func (uc *aggrUseCase) GetTransactions(ctx context.Context, req *GetTransactionsRequest) (*GetTransactionsResponse, error) {
	var (
		cMu  sync.RWMutex
		acMu sync.RWMutex

		categories = make(map[string]*entity.Category)
		accounts   = make(map[string]*entity.Account)
	)

	tRes, err := uc.transactionUseCase.GetTransactions(ctx, req.ToGetTransactionsRequest())
	if err != nil {
		return nil, err
	}
	ts := tRes.Transactions

	tObjs := make([]*TransactionObj, len(ts))
	for i := 0; i < len(ts); i++ {
		tObjs[i] = &TransactionObj{
			Transaction: ts[i],
		}
	}

	if req.CategoryID != nil {
		cRes, err := uc.categoryUseCase.GetCategory(ctx, req.ToGetCategoryRequest())
		if err != nil {
			return nil, err
		}
		categories[req.GetCategoryID()] = cRes.Category
	}

	if req.AccountID != nil {
		acRes, err := uc.accountUseCase.GetAccount(ctx, req.ToGetAccountRequest())
		if err != nil {
			return nil, err
		}
		accounts[req.GetAccountID()] = acRes.Account
	}

	g := new(errgroup.Group)

	// set category
	g.Go(func() error {
		return util.ParallelizeWork(ctx, len(ts), 50, func(ctx context.Context, i int) error {
			t := ts[i]

			// check memory cache for category
			cMu.RLock()
			c, ok := categories[t.GetCategoryID()]
			if ok {
				tObjs[i].Category = c
				cMu.RUnlock()
				return nil
			}
			cMu.RUnlock()

			// if not in memory cache, fetch category from repo and cache it
			cMu.Lock()
			if c == nil {
				cRes, err := uc.categoryUseCase.GetCategory(ctx, &category.GetCategoryRequest{
					UserID:     req.UserID,
					CategoryID: t.CategoryID,
				})
				if err != nil {
					log.Ctx(ctx).Error().Msgf("fail to get category of transaction, transactionID: %v, err: %v",
						t.GetTransactionID(), err)
					return err
				}

				tObjs[i].Category = cRes.Category
				categories[t.GetCategoryID()] = cRes.Category
			}
			cMu.Unlock()

			return nil
		})
	})

	// set account
	g.Go(func() error {
		return util.ParallelizeWork(ctx, len(ts), 50, func(ctx context.Context, i int) error {
			t := ts[i]

			// check memory cache for account
			acMu.RLock()
			ac, ok := accounts[t.GetAccountID()]
			if ok {
				tObjs[i].Account = ac
				acMu.RUnlock()
				return nil
			}
			acMu.RUnlock()

			// if not in memory cache, fetch account from repo and cache it
			acMu.Lock()
			if ac == nil {
				acRes, err := uc.accountUseCase.GetAccount(ctx, &account.GetAccountRequest{
					UserID:    req.UserID,
					AccountID: t.AccountID,
				})
				if err != nil {
					log.Ctx(ctx).Error().Msgf("fail to get account of transaction, transactionID: %v, err: %v",
						t.GetTransactionID(), err)
					return err
				}

				tObjs[i].Account = acRes.Account
				accounts[t.GetAccountID()] = acRes.Account
			}
			acMu.Unlock()

			return nil
		})
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &GetTransactionsResponse{
		TransactionObjs: tObjs,
		Paging:          req.Paging,
	}, nil
}

func (uc *aggrUseCase) GetBudgetWithCategories(
	ctx context.Context,
	req *GetBudgetWithCategoriesRequest,
) (*GetBudgetWithCategoriesResponse, error) {
	budgetRes, err := uc.budgetUseCase.GetBudget(
		ctx,
		&budget.GetBudgetRequest{
			UserID:   req.UserID,
			BudgetID: req.BudgetID,
			Date:     req.Date,
		},
	)
	if err != nil {
		return nil, err
	}

	budget := budgetRes.GetBudget()
	if budget == nil {
		log.Ctx(ctx).Error().Msgf("cannot find budget with budgetID=%s", req.GetBudgetID())
		return nil, fmt.Errorf("cannot find budget with budgetID=%s", req.GetBudgetID())
	}

	categoryRes, err := uc.categoryUseCase.GetCategories(
		ctx,
		&category.GetCategoriesRequest{
			CategoryIDs: budget.GetCategoryIDs(),
		},
	)
	if err != nil {
		return nil, err
	}

	if len(categoryRes.GetCategories()) != len(budget.GetCategoryIDs()) {
		log.Ctx(ctx).Error().Msgf("some categories are missing for ids=%+v", budget.GetCategoryIDs())
		return nil, fmt.Errorf("some categories are missing for ids=%+v", budget.GetCategoryIDs())
	}

	return &GetBudgetWithCategoriesResponse{
		Budget:     budget,
		Categories: categoryRes.GetCategories(),
	}, nil
}

func (uc *aggrUseCase) CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*CreateTransactionResponse, error) {
	tRes, err := uc.transactionUseCase.CreateTransaction(ctx, req.ToCreateTransactionRequest())
	if err != nil {
		return nil, err
	}

	return &CreateTransactionResponse{
		tRes.Transaction,
	}, nil
}

func (uc *aggrUseCase) UpdateTransaction(ctx context.Context, req *UpdateTransactionRequest) (*UpdateTransactionResponse, error) {
	tRes, err := uc.transactionUseCase.UpdateTransaction(ctx, req.ToUpdateTransactionRequest())
	if err != nil {
		return nil, err
	}

	return &UpdateTransactionResponse{
		tRes.Transaction,
	}, nil
}
