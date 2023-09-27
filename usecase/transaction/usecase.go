package transaction

import (
	"context"
	"errors"
	"time"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/common"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidCategoryIDs = errors.New("invalid category_ids")
)

type transactionUseCase struct {
	txMgr            repo.TxMgr
	categoryRepo     repo.CategoryRepo
	accountRepo      repo.AccountRepo
	transactionRepo  repo.TransactionRepo
	budgetRepo       repo.BudgetRepo
	exchangeRateRepo repo.ExchangeRateRepo
}

func NewTransactionUseCase(
	txMgr repo.TxMgr,
	categoryRepo repo.CategoryRepo,
	accountRepo repo.AccountRepo,
	transactionRepo repo.TransactionRepo,
	budgetRepo repo.BudgetRepo,
	exchangeRateRepo repo.ExchangeRateRepo,
) UseCase {
	return &transactionUseCase{
		txMgr,
		categoryRepo,
		accountRepo,
		transactionRepo,
		budgetRepo,
		exchangeRateRepo,
	}
}

func (uc *transactionUseCase) GetTransaction(ctx context.Context, req *GetTransactionRequest) (*GetTransactionResponse, error) {
	t, err := uc.transactionRepo.Get(ctx, req.ToTransactionFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transaction from repo, err: %v", err)
		return nil, err
	}

	c, err := uc.categoryRepo.Get(ctx, req.ToCategoryFilter(t.GetCategoryID()))
	if err != nil && err != repo.ErrCategoryNotFound {
		log.Ctx(ctx).Error().Msgf("fail to get category from repo, err: %v", err)
		return nil, err
	}

	// hide deleted category
	if c != nil {
		t.SetCategory(c)
	} else {
		t.SetCategoryID(goutil.String(""))
	}

	ac, err := uc.accountRepo.Get(ctx, req.ToAccountFilter(t.GetAccountID()))
	if err != nil && err != repo.ErrAccountNotFound {
		log.Ctx(ctx).Error().Msgf("fail to get account from repo, err: %v", err)
		return nil, err
	}

	// deleted account can be shown
	t.SetAccount(ac)

	return &GetTransactionResponse{
		Transaction: t,
	}, nil
}

func (uc *transactionUseCase) GetTransactionGroups(
	ctx context.Context,
	req *GetTransactionGroupsRequest,
) (*GetTransactionGroupsResponse, error) {
	res, err := uc.GetTransactions(ctx, req.GetTransactionsRequest)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transactions, err: %v", err)
		return nil, err
	}

	u := entity.GetUserFromCtx(ctx)

	ers, err := uc.exchangeRateRepo.GetMany(ctx, req.ToExchangeRateFilter(u.Meta.GetCurrency()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get exchange rates from repo, err: %v", err)
		return nil, err
	}

	loc, err := time.LoadLocation(req.AppMeta.GetTimezone())
	if err != nil {
		return nil, entity.ErrInvalidTimezone
	}

	var (
		transactionGroups    = make([]*common.TransactionSummary, 0)
		transactionGroupsMap = make(map[string]*common.TransactionSummary)
	)
	// transactions is already ordered in desc order of transaction_time
	for _, t := range res.Transactions {
		ts := time.UnixMilli(int64(t.GetTransactionTime())).In(loc)
		date := util.FormatDate(ts)

		if _, ok := transactionGroupsMap[date]; !ok {
			ts := &common.TransactionSummary{
				Date:         goutil.String(date),
				Currency:     u.Meta.Currency,
				Sum:          goutil.Float64(0),
				Transactions: make([]*entity.Transaction, 0),
			}
			transactionGroupsMap[date] = ts
			transactionGroups = append(transactionGroups, ts)
		}

		transactionGroup := transactionGroupsMap[date]
		transactionGroup.Transactions = append(transactionGroup.Transactions, t)

		amount := t.GetAmount()
		if u.Meta.GetCurrency() != t.GetCurrency() {
			er := entity.BinarySearchExchangeRates(t, ers)
			amount *= er.GetRate()
		}

		transactionGroup.Sum = goutil.Float64(transactionGroup.GetSum() + amount)
	}

	return &GetTransactionGroupsResponse{
		TransactionGroups: transactionGroups,
		Paging:            req.Paging,
	}, nil
}

func (uc *transactionUseCase) GetTransactions(
	ctx context.Context,
	req *GetTransactionsRequest,
) (*GetTransactionsResponse, error) {
	isDeletedCategory := make(map[string]bool)

	// convert empty category ID to query of deleted categories
	if req.CategoryID != nil && req.GetCategoryID() == "" {
		cs, err := uc.categoryRepo.GetMany(ctx, req.ToCategoryFilter(nil, uint32(entity.CategoryStatusDeleted)))
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get deleted categories from repo, err: %v", err)
			return nil, err
		}

		categoryIDs := make([]string, 0)
		for _, c := range cs {
			categoryIDs = append(categoryIDs, c.GetCategoryID())
			isDeletedCategory[c.GetCategoryID()] = true
		}

		req.CategoryID = nil
		req.CategoryIDs = categoryIDs
	}

	// get transactions
	ts, err := uc.transactionRepo.GetMany(ctx, req.ToTransactionFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transactions from repo, err: %v", err)
		return nil, err
	}

	var (
		categoryIDs = make([]string, 0)
		accountIDs  = make([]string, 0)
	)
	for _, t := range ts {
		accountIDs = append(accountIDs, t.GetAccountID())

		if !isDeletedCategory[t.GetCategoryID()] {
			categoryIDs = append(categoryIDs, t.GetCategoryID())
		}
	}
	categoryIDs = goutil.RemoveDuplicateString(categoryIDs)
	accountIDs = goutil.RemoveDuplicateString(accountIDs)

	// get categories
	var cs []*entity.Category
	if len(categoryIDs) > 0 {
		cs, err = uc.categoryRepo.GetMany(ctx, req.ToCategoryFilter(categoryIDs, uint32(entity.CategoryStatusNormal)))
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get categories from repo, err: %v", err)
			return nil, err
		}
	}

	csMap := make(map[string]*entity.Category)
	for _, c := range cs {
		csMap[c.GetCategoryID()] = c
	}

	// get accounts
	var acs []*entity.Account
	if len(accountIDs) > 0 {
		acs, err = uc.accountRepo.GetMany(ctx, req.ToAccountFilter(accountIDs))
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get accounts from repo, err: %v", err)
			return nil, err
		}
	}

	acsMap := make(map[string]*entity.Account)
	for _, ac := range acs {
		acsMap[ac.GetAccountID()] = ac
	}

	// set accounts and categories
	for _, t := range ts {
		// deleted account can be shown
		ac := acsMap[t.GetAccountID()]
		t.SetAccount(ac)

		// hide deleted category
		if c, ok := csMap[t.GetCategoryID()]; ok {
			t.SetCategory(c)
		} else {
			t.SetCategoryID(goutil.String(""))
		}
	}

	return &GetTransactionsResponse{
		Transactions: ts,
		Paging:       req.Paging,
	}, nil
}

func (uc *transactionUseCase) DeleteTransaction(ctx context.Context, req *DeleteTransactionRequest) (*DeleteTransactionResponse, error) {
	tf := req.ToTransactionFilter()

	t, err := uc.transactionRepo.Get(ctx, tf)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transaction from repo, err: %v", err)
		return nil, err
	}

	if err := uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		if err := uc.transactionRepo.Delete(txCtx, tf); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to mark transaction as deleted, err: %v", err)
			return err
		}

		if t.GetAmount() != 0 {
			// get account
			ac, err := uc.accountRepo.Get(ctx, req.ToAccountFilter(t))
			if err != nil && err != repo.ErrAccountNotFound {
				log.Ctx(txCtx).Error().Msgf("fail to get account from repo, err: %v", err)
				return err
			}

			if err == repo.ErrAccountNotFound {
				log.Ctx(txCtx).Error().Msgf("account has been deleted, account ID: %v", ac.GetAccountID())
				return nil
			}

			// make currency conversion if necessary
			amount, err := uc.getAmountAfterConversion(txCtx, t, ac)
			if err != nil {
				log.Ctx(txCtx).Error().Msgf("fail convert transaction currency, err: %v", err)
				return err
			}

			// reset account balance
			newBalance := ac.GetBalance() - amount
			nac, err := ac.Update(entity.NewAccountUpdate(
				entity.WithUpdateAccountBalance(goutil.Float64(newBalance)),
			))
			if err != nil {
				return err
			}

			if nac != nil {
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

	if err := t.CanTransactionUnderCategory(c); err != nil {
		return nil, err
	}

	ac, err := uc.accountRepo.Get(ctx, req.ToAccountFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get account from repo, err: %v", err)
		return nil, err
	}

	if err := t.CanTransactionUnderAccount(ac); err != nil {
		return nil, err
	}

	if err := uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		// create transaction
		_, err := uc.transactionRepo.Create(txCtx, t)
		if err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to save new transaction to repo, err: %v", err)
			return err
		}

		// make currency conversion if necessary
		amount, err := uc.getAmountAfterConversion(txCtx, t, ac)
		if err != nil {
			log.Ctx(txCtx).Error().Msgf("fail convert transaction currency, err: %v", err)
			return err
		}

		// update account balance
		newBalance := ac.GetBalance() + amount
		nac, err := ac.Update(entity.NewAccountUpdate(
			entity.WithUpdateAccountBalance(goutil.Float64(newBalance)),
		))
		if err != nil {
			return err
		}

		if nac != nil {
			if err := uc.accountRepo.Update(txCtx, req.ToAccountFilter(), nac); err != nil {
				log.Ctx(txCtx).Error().Msgf("fail to update account balance, err: %v", err)
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	t.SetCategory(c)
	t.SetAccount(ac)

	return &CreateTransactionResponse{
		Transaction: t,
	}, nil
}

func (uc *transactionUseCase) UpdateTransaction(ctx context.Context, req *UpdateTransactionRequest) (*UpdateTransactionResponse, error) {
	t, err := uc.transactionRepo.Get(ctx, req.ToTransactionFilter())
	if err != nil {
		return nil, err
	}
	oldT := *t

	tu := t.Update(req.ToTransactionUpdate())
	if tu == nil {
		log.Ctx(ctx).Info().Msg("transaction has no updates")
		return &UpdateTransactionResponse{
			Transaction: t,
		}, nil
	}

	oldAc, err := uc.accountRepo.Get(ctx, req.ToAccountFilter(oldT.GetAccountID()))
	if err != nil && err != repo.ErrAccountNotFound {
		log.Ctx(ctx).Info().Msgf("fail to get old account from repo, err: %v", err)
		return nil, err
	}

	newAc := oldAc
	if tu.AccountID != nil {
		newAc, err = uc.accountRepo.Get(ctx, req.ToAccountFilter(tu.GetAccountID()))
		if err != nil {
			log.Ctx(ctx).Info().Msgf("fail to get new account from repo, err: %v", err)
			return nil, err
		}

		if err := t.CanTransactionUnderAccount(newAc); err != nil {
			return nil, err
		}
	}

	if tu.CategoryID != nil {
		newCategory, err := uc.categoryRepo.Get(ctx, req.ToCategoryFilter())
		if err != nil {
			log.Ctx(ctx).Info().Msgf("fail to get new category from repo, err: %v", err)
			return nil, err
		}

		if err := t.CanTransactionUnderCategory(newCategory); err != nil {
			return nil, err
		}
	}

	if err = uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		// save updates
		if err = uc.transactionRepo.Update(txCtx, req.ToTransactionFilter(), tu); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to save transaction updates to repo, err: %v", err)
			return err
		}

		if tu.AccountID != nil || tu.Amount != nil {
			// revert balance of old account
			if oldAc != nil {
				amount, err := uc.getAmountAfterConversion(txCtx, &oldT, oldAc)
				if err != nil {
					log.Ctx(txCtx).Error().Msgf("fail convert transaction currency, err: %v", err)
					return err
				}

				oldAcBalance := oldAc.GetBalance() - amount
				acu, err := oldAc.Update(entity.NewAccountUpdate(
					entity.WithUpdateAccountBalance(goutil.Float64(oldAcBalance)),
				))
				if err != nil {
					return err
				}

				if acu != nil {
					if err := uc.accountRepo.Update(txCtx, req.ToAccountFilter(oldT.GetAccountID()), acu); err != nil {
						log.Ctx(txCtx).Error().Msgf("fail to update old account balance, err: %v", err)
						return err
					}
				}
			}

			// update balance of new account
			if newAc != nil {
				amount, err := uc.getAmountAfterConversion(txCtx, t, newAc)
				if err != nil {
					log.Ctx(txCtx).Error().Msgf("fail convert transaction currency, err: %v", err)
					return err
				}

				newAcBalance := newAc.GetBalance() + amount
				acu, err := newAc.Update(entity.NewAccountUpdate(
					entity.WithUpdateAccountBalance(goutil.Float64(newAcBalance)),
				))
				if err != nil {
					return err
				}

				if acu != nil {
					if err := uc.accountRepo.Update(txCtx, req.ToAccountFilter(newAc.GetAccountID()), acu); err != nil {
						log.Ctx(txCtx).Error().Msgf("fail to update new account balance, err: %v", err)
						return err
					}
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

func (uc *transactionUseCase) SumTransactions(ctx context.Context, req *SumTransactionsRequest) (*SumTransactionsResponse, error) {
	ts, err := uc.transactionRepo.GetMany(ctx, req.ToTransactionFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transactions from repo, err: %v", err)
		return nil, err
	}

	u := entity.GetUserFromCtx(ctx)

	ers, err := uc.exchangeRateRepo.GetMany(ctx, req.ToExchangeRateFilter(u.Meta.GetCurrency()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get exchange rates from repo, err: %v", err)
		return nil, err
	}

	sumByTT := make(map[uint32]float64)
	for tt := range entity.TransactionTypes {
		if (req.TransactionType == nil) ||
			(req.TransactionType != nil && tt == req.GetTransactionType()) {
			sumByTT[tt] = 0
		}
	}

	for _, t := range ts {
		amount := t.GetAmount()

		if t.GetCurrency() != u.Meta.GetCurrency() {
			er := entity.BinarySearchExchangeRates(t, ers)
			amount *= er.GetRate()
		}

		sumByTT[t.GetTransactionType()] += amount
	}

	sums := make([]*common.TransactionSummary, 0)
	for tt, sum := range sumByTT {
		sums = append(sums, &common.TransactionSummary{
			TransactionType: goutil.Uint32(tt),
			Sum:             goutil.Float64(util.RoundFloatToStandardDP(sum)),
			Currency:        u.Meta.Currency,
		})
	}

	return &SumTransactionsResponse{
		Sums: sums,
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

func (uc *transactionUseCase) getAmountAfterConversion(ctx context.Context, t *entity.Transaction, ac *entity.Account) (float64, error) {
	amount := t.GetAmount()

	if t.GetCurrency() == ac.GetCurrency() {
		return amount, nil
	}

	er, err := uc.exchangeRateRepo.Get(ctx, &repo.GetExchangeRateFilter{
		Timestamp: t.TransactionTime,
		From:      t.Currency,
		To:        ac.Currency,
	})
	if err != nil {
		return 0, err
	}
	amount *= er.GetRate()

	return amount, nil
}
