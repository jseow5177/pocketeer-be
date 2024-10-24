package transaction

import (
	"context"
	"errors"
	"math"
	"sort"
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

// intermediate state
type transactionSummary struct {
	Date         string
	Sum          float64
	TotalExpense float64
	TotalIncome  float64
	Transactions []*entity.Transaction
}

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

	if t.IsTransfer() {
		fromAc, err := uc.accountRepo.Get(ctx, req.ToAccountFilter(t.GetFromAccountID()))
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get from_account from repo, err: %v", err)
			return nil, err
		}

		t.SetFromAccount(fromAc)

		toAc, err := uc.accountRepo.Get(ctx, req.ToAccountFilter(t.GetToAccountID()))
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get to_account from repo, err: %v", err)
			return nil, err
		}

		t.SetToAccount(toAc)
	} else {
		c, err := uc.categoryRepo.Get(ctx, req.ToCategoryFilter(t.GetCategoryID()))
		if err != nil && err != repo.ErrCategoryNotFound {
			log.Ctx(ctx).Error().Msgf("fail to get category from repo, err: %v", err)
			return nil, err
		}

		// hide deleted or empty category category
		if c == nil || c.IsDeleted() {
			t.SetCategoryID(goutil.String(""))
		} else {
			t.SetCategory(c)
		}

		ac, err := uc.accountRepo.Get(ctx, req.ToAccountFilter(t.GetAccountID()))
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get account from repo, err: %v", err)
			return nil, err
		}

		t.SetAccount(ac)
	}

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

	loc, err := time.LoadLocation(req.AppMeta.GetTimezone())
	if err != nil {
		return nil, entity.ErrInvalidTimezone
	}

	u := entity.GetUserFromCtx(ctx)

	var (
		transactionGroups    = make([]*transactionSummary, 0)
		transactionGroupsMap = make(map[string]*transactionSummary)
	)
	// transactions is already ordered in desc order of transaction_time
	for _, t := range res.Transactions {
		ts := time.UnixMilli(int64(t.GetTransactionTime())).In(loc)
		date := util.FormatDate(ts)

		if _, ok := transactionGroupsMap[date]; !ok {
			ts := &transactionSummary{
				Date:         date,
				Sum:          0,
				TotalExpense: 0,
				TotalIncome:  0,
				Transactions: make([]*entity.Transaction, 0),
			}
			transactionGroupsMap[date] = ts
			transactionGroups = append(transactionGroups, ts)
		}

		transactionGroup := transactionGroupsMap[date]
		transactionGroup.Transactions = append(transactionGroup.Transactions, t)

		var amount float64
		if !t.IsTransfer() {
			amount, err = uc.getAmountAfterConversion(ctx, t, u.Meta.GetCurrency())
			if err != nil {
				log.Ctx(ctx).Error().Msgf("fail convert transaction currency, err: %v", err)
				return nil, err
			}
		}

		if t.IsExpense() {
			transactionGroup.TotalExpense += amount
		} else if t.IsIncome() {
			transactionGroup.TotalIncome += amount
		}

		transactionGroup.Sum += amount
	}

	// round after sum
	tgs := make([]*common.Summary, 0)
	for _, transactionGroup := range transactionGroups {
		tgs = append(tgs, common.NewSummary(
			common.WithSummaryCurrency(u.Meta.Currency),
			common.WithSummaryDate(goutil.String(transactionGroup.Date)),
			common.WithSummaryTransactions(transactionGroup.Transactions),
			common.WithSummarySum(goutil.Float64(transactionGroup.Sum)),
			common.WithSummaryTotalExpense(goutil.Float64(transactionGroup.TotalExpense)),
			common.WithSummaryTotalIncome(goutil.Float64(transactionGroup.TotalIncome)),
		))
	}

	return &GetTransactionGroupsResponse{
		TransactionGroups: tgs,
		Paging:            req.Paging,
	}, nil
}

func (uc *transactionUseCase) GetTransactions(
	ctx context.Context,
	req *GetTransactionsRequest,
) (*GetTransactionsResponse, error) {
	// convert empty category ID to query of deleted categories
	if req.CategoryID != nil && req.GetCategoryID() == "" {
		cf := req.ToCategoryFilter(
			nil,
			uint32(entity.CategoryStatusDeleted),
		)

		cs, err := uc.categoryRepo.GetMany(ctx, cf)
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get deleted categories from repo, err: %v", err)
			return nil, err
		}

		categoryIDs := make([]string, 0)
		for _, c := range cs {
			categoryIDs = append(categoryIDs, c.GetCategoryID())
		}

		req.CategoryID = nil
		req.CategoryIDs = categoryIDs
	}

	// get transactions
	ts, err := uc.transactionRepo.GetMany(ctx, req.ToTransactionQuery())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transactions from repo, err: %v", err)
		return nil, err
	}

	var (
		categoryIDs = make([]string, 0)
		accountIDs  = make([]string, 0)
	)
	for _, t := range ts {
		if t.IsTransfer() {
			accountIDs = append(accountIDs, t.GetFromAccountID(), t.GetToAccountID())
		} else {
			accountIDs = append(accountIDs, t.GetAccountID())
		}

		if t.GetCategoryID() != "" {
			categoryIDs = append(categoryIDs, t.GetCategoryID())
		}
	}

	categoryIDs = goutil.RemoveDuplicateString(categoryIDs)
	accountIDs = goutil.RemoveDuplicateString(accountIDs)

	// get categories
	var cs []*entity.Category
	if len(categoryIDs) > 0 {
		cf := req.ToCategoryFilter(
			categoryIDs,
			uint32(entity.CategoryStatusNormal),
		)

		cs, err = uc.categoryRepo.GetMany(ctx, cf)
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
		if t.IsTransfer() {
			fromAc := acsMap[t.GetFromAccountID()]
			toAc := acsMap[t.GetToAccountID()]

			t.SetFromAccount(fromAc)
			t.SetToAccount(toAc)
		} else {
			ac := acsMap[t.GetAccountID()]
			t.SetAccount(ac)
		}

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
		tu, err := t.Update(
			entity.WithUpdateTransactionStatus(goutil.Uint32(uint32(entity.TransactionStatusDeleted))),
		)
		if err != nil {
			return err
		}

		// mark transaction as deleted
		if err := uc.transactionRepo.Update(txCtx, tf, tu); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to save transaction updates to repo, err: %v", err)
			return err
		}

		if t.GetAmount() == 0 {
			return nil
		}

		// left: to minus amount, right: to plus amount
		acs := make(map[*entity.Account]*entity.Account)

		if t.IsTransfer() {
			fromAc, err := uc.accountRepo.Get(txCtx, req.ToAccountFilter(t.GetFromAccountID()))
			if err != nil && err != repo.ErrAccountNotFound {
				log.Ctx(txCtx).Error().Msgf("fail to get from_account from repo, err: %v", err)
				return err
			}

			toAc, err := uc.accountRepo.Get(txCtx, req.ToAccountFilter(t.GetToAccountID()))
			if err != nil && err != repo.ErrAccountNotFound {
				log.Ctx(txCtx).Error().Msgf("fail to get to_account from repo, err: %v", err)
				return err
			}

			acs[toAc] = fromAc
		} else {
			ac, err := uc.accountRepo.Get(txCtx, req.ToAccountFilter(t.GetAccountID()))
			if err != nil && err != repo.ErrAccountNotFound {
				log.Ctx(txCtx).Error().Msgf("fail to get account from repo, err: %v", err)
				return err
			}

			acs[ac] = nil
		}

		for minusAc, addAc := range acs {
			if minusAc != nil {
				if err := uc.updateAccountBalance(txCtx, t, minusAc, false); err != nil {
					log.Ctx(txCtx).Error().Msgf("fail to update account balance, err: %v", err)
					return err
				}
			}

			if addAc != nil {
				if err := uc.updateAccountBalance(txCtx, t, addAc, true); err != nil {
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
	t, err := req.ToTransactionEntity()
	if err != nil {
		return nil, err
	}

	// left: to minus amount, right: to plus amount
	acs := make(map[*entity.Account]*entity.Account)

	if t.IsTransfer() {
		fromAc, err := uc.accountRepo.Get(ctx, req.ToAccountFilter(req.GetFromAccountID()))
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get from_account from repo, account_id: %v, err: %v",
				req.GetFromAccountID(), err)
			return nil, err
		}
		t.SetFromAccount(fromAc)

		toAc, err := uc.accountRepo.Get(ctx, req.ToAccountFilter(req.GetToAccountID()))
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get to_account from repo, account_id: %v, err: %v",
				req.GetToAccountID(), err)
			return nil, err
		}
		t.SetToAccount(toAc)

		acs[fromAc] = toAc
	} else {
		c, err := uc.categoryRepo.Get(ctx, req.ToCategoryFilter())
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get category from repo, category_id: %v, err: %v",
				req.GetCategoryID(), err)
			return nil, err
		}
		t.SetCategory(c)

		if err := t.CanTransactionUnderCategory(c); err != nil {
			return nil, err
		}

		ac, err := uc.accountRepo.Get(ctx, req.ToAccountFilter(req.GetAccountID()))
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get account from repo, account_id: %v, err: %v",
				req.GetAccountID(), err)
			return nil, err
		}

		if err := t.CanTransactionUnderAccount(ac); err != nil {
			return nil, err
		}
		t.SetAccount(ac)

		acs[nil] = ac
	}

	if err := uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		// create transaction
		_, err := uc.transactionRepo.Create(txCtx, t)
		if err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to save new transaction to repo, err: %v", err)
			return err
		}

		for minusAc, addAc := range acs {
			if minusAc != nil {
				if err := uc.updateAccountBalance(txCtx, t, minusAc, false); err != nil {
					log.Ctx(txCtx).Error().Msgf("fail to update account balance, err: %v", err)
					return err
				}
			}

			if addAc != nil {
				if err := uc.updateAccountBalance(txCtx, t, addAc, true); err != nil {
					log.Ctx(txCtx).Error().Msgf("fail to update account balance, err: %v", err)
					return err
				}
			}
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

	oldT, err := t.Clone()
	if err != nil {
		return nil, err
	}

	tu, err := t.Update(
		entity.WithUpdateTransactionAmount(req.Amount),
		entity.WithUpdateTransactionTime(req.TransactionTime),
		entity.WithUpdateTransactionNote(req.Note),
		entity.WithUpdateTransactionAccountID(req.AccountID),
		entity.WithUpdateTransactionCategoryID(req.CategoryID),
		entity.WithUpdateTransactionFromAccountID(req.FromAccountID),
		entity.WithUpdateTransactionToAccountID(req.ToAccountID),
		entity.WithUpdateTransactionType(req.TransactionType),
		entity.WithUpdateTransactionCurrency(req.Currency),
	)
	if err != nil {
		return nil, err
	}

	if tu == nil {
		log.Ctx(ctx).Info().Msg("transaction has no updates")
		return &UpdateTransactionResponse{
			Transaction: t,
		}, nil
	}

	var (
		needResetAccountBalance = tu.TransactionTime != nil || tu.Currency != nil ||
			tu.TransactionType != nil || tu.Amount != nil || tu.AccountID != nil ||
			tu.FromAccountID != nil || tu.ToAccountID != nil

		transactionOps = []*struct {
			AccountID string
			*entity.Transaction
			Add bool
		}{
			{AccountID: oldT.GetAccountID(), Add: false, Transaction: oldT},
			{AccountID: t.GetAccountID(), Add: true, Transaction: t},

			{AccountID: oldT.GetFromAccountID(), Add: true, Transaction: oldT},
			{AccountID: t.GetFromAccountID(), Add: false, Transaction: t},

			{AccountID: oldT.GetToAccountID(), Add: false, Transaction: oldT},
			{AccountID: t.GetToAccountID(), Add: true, Transaction: t},
		}

		accountIDs = []string{oldT.GetAccountID(), t.GetAccountID(),
			oldT.GetFromAccountID(), t.GetFromAccountID(), oldT.GetToAccountID(), t.GetToAccountID()}
	)

	acs := make(map[string]*entity.Account)
	if needResetAccountBalance {
		acf := repo.NewAccountFilter(
			req.GetUserID(),
			repo.WithAccountIDs(goutil.RemoveDuplicateString(accountIDs)),
		)

		accounts, err := uc.accountRepo.GetMany(ctx, acf)
		if err != nil {
			log.Ctx(ctx).Info().Msgf("fail to get accounts from repo, err: %v", err)
			return nil, err
		}

		for _, ac := range accounts {
			acs[ac.GetAccountID()] = ac
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
		if err := uc.transactionRepo.Update(txCtx, req.ToTransactionFilter(), tu); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to save transaction updates to repo, err: %v", err)
			return err
		}

		for _, op := range transactionOps {
			ac, ok := acs[op.AccountID]
			if !ok || ac == nil {
				continue
			}

			if err := uc.updateAccountBalance(txCtx, op.Transaction, ac, op.Add); err != nil {
				log.Ctx(txCtx).Error().Msgf("fail to update account balance, err: %v", err)
				return err
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
	ts, err := uc.transactionRepo.GetMany(ctx, req.ToTransactionQuery())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transactions from repo, err: %v", err)
		return nil, err
	}

	sumByTT := make(map[uint32]float64)
	for tt := range entity.TransactionTypes {
		if (req.TransactionType == nil) ||
			(req.TransactionType != nil && tt == req.GetTransactionType()) {
			sumByTT[tt] = 0
		}
	}

	u := entity.GetUserFromCtx(ctx)

	for _, t := range ts {
		amount, err := uc.getAmountAfterConversion(ctx, t, u.Meta.GetCurrency())
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail convert transaction currency, err: %v", err)
			return nil, err
		}

		sumByTT[t.GetTransactionType()] += amount
	}

	sums := make([]*common.Summary, 0)
	for tt, sum := range sumByTT {
		sums = append(sums, common.NewSummary(
			common.WithSummaryTransactionType(goutil.Uint32(tt)),
			common.WithSummarySum(goutil.Float64(sum)),
			common.WithSummaryCurrency(u.Meta.Currency),
		))
	}

	return &SumTransactionsResponse{
		Sums: sums,
	}, nil
}

func (uc *transactionUseCase) GetTransactionsSummary(ctx context.Context, req *GetTransactionsSummaryRequest) (*GetTransactionsSummaryResponse, error) {
	user := req.GetUser()

	tq, err := req.ToTransactionQuery()
	if err != nil {
		return nil, err
	}

	// get latest transactions
	transactions, err := uc.transactionRepo.GetMany(ctx, tq)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transactions from repo, err: %v", err)
		return nil, err
	}

	loc, err := time.LoadLocation(req.AppMeta.GetTimezone())
	if err != nil {
		return nil, entity.ErrInvalidTimezone
	}

	var (
		now                 = time.Now().In(loc)
		date                time.Time
		monthDiff, yearDiff int
	)

	switch req.GetUnit() {
	case uint32(entity.SnapshotUnitMonth):
		date = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()) // start of month
		monthDiff = -1
	default:
		return nil, entity.ErrInvalidSnapshotUnit
	}

	var (
		i = req.GetInterval()
		d = util.FormatDate(date)

		transactionGroupsMap = map[string]*transactionSummary{
			d: {
				Date:         d,
				Sum:          0,
				TotalExpense: 0,
				TotalIncome:  0,
			},
		}
	)
	for i > 0 {
		date = date.AddDate(yearDiff, monthDiff, 0)
		d := util.FormatDate(date)

		transactionGroupsMap[d] = &transactionSummary{
			Date:         d,
			Sum:          0,
			TotalExpense: 0,
			TotalIncome:  0,
		}

		i--
	}

	for _, transaction := range transactions {
		var (
			date time.Time
			t    = time.UnixMilli(int64(transaction.GetTransactionTime())).In(loc)
		)

		switch req.GetUnit() {
		case uint32(entity.SnapshotUnitMonth):
			date = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()) // start of month
		}

		d := util.FormatDate(date)

		// shouldn't happen
		if _, ok := transactionGroupsMap[d]; !ok {
			transactionGroupsMap[d] = &transactionSummary{
				Date:         d,
				Sum:          0,
				TotalExpense: 0,
				TotalIncome:  0,
			}
		}

		var amount float64
		if !transaction.IsTransfer() {
			amount, err = uc.getAmountAfterConversion(ctx, transaction, user.Meta.GetCurrency())
			if err != nil {
				log.Ctx(ctx).Error().Msgf("fail convert transaction currency, err: %v", err)
				return nil, err
			}
		}

		transactionGroup := transactionGroupsMap[d]

		if transaction.IsExpense() {
			transactionGroup.TotalExpense += amount
		} else if transaction.IsIncome() {
			transactionGroup.TotalIncome += amount
		}

		transactionGroup.Sum += amount
	}

	dates := make([]string, 0)
	for date := range transactionGroupsMap {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	resp := new(GetTransactionsSummaryResponse)
	for i, date := range dates {
		tg := transactionGroupsMap[date]

		var (
			savings      = tg.Sum
			totalExpense = tg.TotalExpense
			totalIncome  = tg.TotalIncome

			percentSavingsChange, percentTotalExpenseChange, percentTotalIncomeChange    *float64
			absoluteSavingsChange, absoluteTotalExpenseChange, absoluteTotalIncomeChange *float64
		)
		// only calculate change between first and last
		if i > 0 && i == len(dates)-1 {
			var (
				firstSavings      = resp.Savings[0].GetSum()
				firstTotalExpense = resp.TotalExpense[0].GetSum()
				firstTotalIncome  = resp.TotalIncome[0].GetSum()
			)
			absoluteSavingsChange = goutil.Float64(savings - firstSavings)
			if *absoluteSavingsChange == 0 {
				percentSavingsChange = goutil.Float64(0)
			} else if firstSavings != 0 {
				percentSavingsChange = goutil.Float64(*absoluteSavingsChange * 100 / math.Abs(firstSavings))
			}

			absoluteTotalExpenseChange = goutil.Float64(totalExpense - firstTotalExpense)
			if *absoluteTotalExpenseChange == 0 {
				percentTotalExpenseChange = goutil.Float64(0)
			} else if firstTotalExpense != 0 {
				percentTotalExpenseChange = goutil.Float64(*absoluteTotalExpenseChange * 100 / math.Abs(firstTotalExpense))
			}

			absoluteTotalIncomeChange = goutil.Float64(totalIncome - firstTotalIncome)
			if *absoluteTotalIncomeChange == 0 {
				percentTotalIncomeChange = goutil.Float64(0)
			} else if firstTotalIncome != 0 {
				percentTotalIncomeChange = goutil.Float64(*absoluteTotalIncomeChange * 100 / math.Abs(firstTotalIncome))
			}
		}

		// savings
		resp.Savings = append(resp.Savings, common.NewSummary(
			common.WithSummaryDate(goutil.String(date)),
			common.WithSummarySum(goutil.Float64(savings)),
			common.WithSummaryCurrency(user.Meta.Currency),
			common.WithSummaryPercentChange(percentSavingsChange),
			common.WithSummaryAbsoluteChange(absoluteSavingsChange),
		))

		// total expense
		resp.TotalExpense = append(resp.TotalExpense, common.NewSummary(
			common.WithSummaryDate(goutil.String(date)),
			common.WithSummarySum(goutil.Float64(totalExpense)),
			common.WithSummaryCurrency(user.Meta.Currency),
			common.WithSummaryPercentChange(percentTotalExpenseChange),
			common.WithSummaryAbsoluteChange(absoluteTotalExpenseChange),
		))

		// total income
		resp.TotalIncome = append(resp.TotalIncome, common.NewSummary(
			common.WithSummaryDate(goutil.String(date)),
			common.WithSummarySum(goutil.Float64(totalIncome)),
			common.WithSummaryCurrency(user.Meta.Currency),
			common.WithSummaryPercentChange(percentTotalIncomeChange),
			common.WithSummaryAbsoluteChange(absoluteTotalIncomeChange),
		))
	}

	return resp, nil
}

func (uc *transactionUseCase) getAmountAfterConversion(ctx context.Context, t *entity.Transaction, currency string) (float64, error) {
	amount := t.GetAmount()

	if t.GetCurrency() == currency {
		return amount, nil
	}

	erf := repo.NewExchangeRateFilter(
		repo.WithExchangeRateFrom(t.Currency),
		repo.WithExchangeRateTo(goutil.String(currency)),
		repo.WithExchangeRateTimestamp(t.TransactionTime),
	)

	er, err := uc.exchangeRateRepo.Get(ctx, erf)
	if err != nil {
		return 0, err
	}

	amount *= er.GetRate()

	return amount, nil
}

func (uc *transactionUseCase) updateAccountBalance(ctx context.Context, t *entity.Transaction, ac *entity.Account, add bool) error {
	// make currency conversion if necessary
	amount, err := uc.getAmountAfterConversion(ctx, t, ac.GetCurrency())
	if err != nil {
		return err
	}

	var newBalance float64
	if add {
		newBalance = ac.GetBalance() + amount
	} else {
		newBalance = ac.GetBalance() - amount
	}

	nac, err := ac.Update(
		entity.WithUpdateAccountBalance(goutil.Float64(newBalance)),
	)
	if err != nil {
		return err
	}

	if nac != nil {
		if err := uc.accountRepo.Update(ctx, repo.NewAccountFilter(
			ac.GetUserID(),
			repo.WithAccountID(ac.AccountID),
		), nac); err != nil {
			return err
		}
	}

	return nil
}
