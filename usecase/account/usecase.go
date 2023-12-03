package account

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/common"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

type accountUseCase struct {
	txMgr            repo.TxMgr
	accountRepo      repo.AccountRepo
	transactionRepo  repo.TransactionRepo
	holdingRepo      repo.HoldingRepo
	lotRepo          repo.LotRepo
	quoteRepo        repo.QuoteRepo
	securityRepo     repo.SecurityRepo
	exchangeRateRepo repo.ExchangeRateRepo
	snapshotRepo     repo.SnapshotRepo
}

func NewAccountUseCase(
	txMgr repo.TxMgr,
	accountRepo repo.AccountRepo,
	transactionRepo repo.TransactionRepo,
	holdingRepo repo.HoldingRepo,
	lotRepo repo.LotRepo,
	quoteRepo repo.QuoteRepo,
	securityRepo repo.SecurityRepo,
	exchangeRateRepo repo.ExchangeRateRepo,
	snapshotRepo repo.SnapshotRepo,
) UseCase {
	return &accountUseCase{
		txMgr,
		accountRepo,
		transactionRepo,
		holdingRepo,
		lotRepo,
		quoteRepo,
		securityRepo,
		exchangeRateRepo,
		snapshotRepo,
	}
}

func (uc *accountUseCase) GetAccount(ctx context.Context, req *GetAccountRequest) (*GetAccountResponse, error) {
	ac, err := uc.accountRepo.Get(ctx, req.ToAccountFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get account from repo, err: %v", err)
		return nil, err
	}

	if !ac.IsInvestment() {
		return &GetAccountResponse{
			Account: ac,
		}, nil
	}

	if err := uc.getAccountHoldingsAndLots(ctx, ac); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get account holdings and lots, err: %v", err)
		return nil, err
	}

	return &GetAccountResponse{
		Account: ac,
	}, nil
}

func (uc *accountUseCase) GetAccounts(ctx context.Context, req *GetAccountsRequest) (*GetAccountsResponse, error) {
	acs, err := uc.accountRepo.GetMany(ctx, req.ToAccountFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get accounts from repo, err: %v", err)
		return nil, err
	}

	var (
		mu  sync.Mutex
		now = time.Now().UnixMilli()
		u   = entity.GetUserFromCtx(ctx)

		assetValue, debtValue float64
	)
	if err := goutil.ParallelizeWork(ctx, len(acs), 10, func(ctx context.Context, workNum int) error {
		ac := acs[workNum]

		if ac.IsInvestment() {
			if err := uc.getAccountHoldingsAndLots(ctx, ac); err != nil {
				log.Ctx(ctx).Error().Msgf("fail to get account holdings and lots, err: %v", err)
				return err
			}
		}

		balance := ac.GetBalance()
		if u.Meta.GetCurrency() != ac.GetCurrency() {
			erf := repo.NewExchangeRateFilter(
				repo.WithExchangeRateFrom(ac.Currency),
				repo.WithExchangeRateTo(u.Meta.Currency),
				repo.WithExchangeRateTimestamp(goutil.Uint64(uint64(now))),
			)
			er, err := uc.exchangeRateRepo.Get(ctx, erf)
			if err != nil {
				log.Ctx(ctx).Error().Msgf("fail to get exchange rate from repo, err: %v", err)
				return err
			}

			balance *= er.GetRate()
		}

		mu.Lock()
		if ac.IsAsset() {
			assetValue += balance
		} else if ac.IsDebt() {
			debtValue += balance
		}
		mu.Unlock()

		return nil
	}); err != nil {
		return nil, err
	}

	netWorth := assetValue + debtValue

	return &GetAccountsResponse{
		NetWorth:   goutil.Float64(util.RoundFloatToStandardDP(netWorth)),
		AssetValue: goutil.Float64(util.RoundFloatToStandardDP(assetValue)),
		DebtValue:  goutil.Float64(util.RoundFloatToStandardDP(debtValue)),
		Currency:   u.Meta.Currency,
		Accounts:   acs,
	}, nil
}

func (uc *accountUseCase) CreateAccount(ctx context.Context, req *CreateAccountRequest) (*CreateAccountResponse, error) {
	ac, err := req.ToAccountEntity()
	if err != nil {
		return nil, err
	}

	if _, err := uc.accountRepo.Create(ctx, ac); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save new account to repo, err: %v", err)
		return nil, err
	}

	if err := uc.computeAccountCostGainAndBalance(ctx, ac); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to compute account cost, gain, and balance, err: %v", err)
		return nil, err
	}

	return &CreateAccountResponse{
		Account: ac,
	}, nil
}

func (uc *accountUseCase) UpdateAccount(ctx context.Context, req *UpdateAccountRequest) (*UpdateAccountResponse, error) {
	ac, err := uc.accountRepo.Get(ctx, req.ToAccountFilter())
	if err != nil {
		return nil, err
	}
	oldBalance := ac.GetBalance()

	acu, err := ac.Update(
		entity.WithUpdateAccountBalance(req.Balance),
		entity.WithUpdateAccountName(req.AccountName),
		entity.WithUpdateAccountNote(req.Note),
	)
	if err != nil {
		return nil, err
	}

	if acu == nil {
		log.Ctx(ctx).Info().Msg("acount has no updates")
		return &UpdateAccountResponse{
			Account: ac,
		}, nil
	}

	if err := uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		if err := uc.accountRepo.Update(txCtx, req.ToAccountFilter(), acu); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to save account updates to repo, err: %v", err)
			return err
		}

		if acu.Balance != nil && req.NeedOffsetTransaction() {
			balanceChange := acu.GetBalance() - oldBalance

			t, err := uc.newUnrecordedTransaction(ac, balanceChange)
			if err != nil {
				log.Ctx(txCtx).Error().Msgf("fail to new unrecorded transaction, err: %v", err)
				return err
			}

			if _, err := uc.transactionRepo.Create(txCtx, t); err != nil {
				log.Ctx(txCtx).Error().Msgf("fail to create unrecorded transaction, err: %v", err)
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &UpdateAccountResponse{
		Account: ac,
	}, nil
}

func (uc *accountUseCase) DeleteAccount(ctx context.Context, req *DeleteAccountRequest) (*DeleteAccountResponse, error) {
	acf := req.ToAccountFilter()

	ac, err := uc.accountRepo.Get(ctx, acf)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get account from repo, err: %v", err)
		return nil, err
	}

	if err := uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		acu, err := ac.Update(
			entity.WithUpdateAccountStatus(goutil.Uint32(uint32(entity.AccountStatusDeleted))),
		)
		if err != nil {
			return err
		}

		if err = uc.accountRepo.Update(txCtx, acf, acu); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to mark account as deleted, err: %v", err)
			return err
		}

		if !ac.IsInvestment() {
			return nil
		}

		hs, err := uc.holdingRepo.GetMany(txCtx, req.ToHoldingFilter())
		if err != nil {
			log.Ctx(txCtx).Error().Msgf("fail get holdings from repo, err: %v", err)
			return err
		}

		holdingIDs := make([]string, 0)
		for _, h := range hs {
			holdingIDs = append(holdingIDs, h.GetHoldingID())
		}

		// mark holdings as deleted
		hu := &entity.HoldingUpdate{
			HoldingStatus: goutil.Uint32(uint32(entity.HoldingStatusDeleted)),
		}
		if err := uc.holdingRepo.UpdateMany(ctx, req.ToHoldingFilter(), hu); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to mark account holdings as deleted, err: %v", err)
			return err
		}

		// mark lots as deleted
		lu := &entity.LotUpdate{
			LotStatus: goutil.Uint32(uint32(entity.LotStatusDeleted)),
		}
		if err := uc.lotRepo.UpdateMany(ctx, req.ToLotFilter(holdingIDs), lu); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to mark lots as deleted, err: %v", err)
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return new(DeleteAccountResponse), nil
}

func (uc *accountUseCase) newUnrecordedTransaction(account *entity.Account, amount float64) (*entity.Transaction, error) {
	tt := uint32(entity.GetTransactionTypeByAmount(amount))

	note := fmt.Sprintf("Unrecorded %s", entity.TransactionTypes[tt])

	return entity.NewTransaction(
		account.GetUserID(),
		entity.WithTransactionAccountID(account.AccountID),
		entity.WithTransactionAmount(goutil.Float64(amount)),
		entity.WithTransactionCurrency(account.Currency),
		entity.WithTransactionType(goutil.Uint32(tt)),
		entity.WithTransactionNote(goutil.String(note)),
		entity.WithTransactionTime(goutil.Uint64(uint64(time.Now().UnixMilli()))),
	)
}

func (uc *accountUseCase) getAccountHoldingsAndLots(ctx context.Context, ac *entity.Account) error {
	hs, err := uc.holdingRepo.GetMany(ctx, repo.NewHoldingFilter(
		repo.WithHoldingUserID(ac.UserID),
		repo.WithHoldingAccountID(ac.AccountID),
	))
	if err != nil {
		return fmt.Errorf("fail to get holdings from repo, err: %v", err)
	}

	for _, h := range hs {
		if h.IsDefault() {
			q, err := uc.quoteRepo.Get(ctx, repo.NewQuoteFilter(
				repo.WithQuoteSymbol(h.Symbol),
			))
			if err != nil {
				return fmt.Errorf("fail to get quote from repo, err: %v", err)
			}
			h.SetQuote(q)

			ls, err := uc.lotRepo.GetMany(ctx, repo.NewLotFilter(
				ac.GetUserID(),
				repo.WithLotHoldingID(h.HoldingID),
			))
			if err != nil {
				return fmt.Errorf("fail to get lots from repo, err: %v", err)
			}
			h.SetLots(ls)
		}

		h.ComputeCostGainAndValue()
	}

	ac.SetHoldings(hs)

	if err := uc.computeAccountCostGainAndBalance(ctx, ac); err != nil {
		return fmt.Errorf("fail to compute account cost, gain, and balance, err: %v", err)
	}

	return nil
}

func (uc *accountUseCase) computeAccountCostGainAndBalance(ctx context.Context, ac *entity.Account) error {
	if !ac.IsInvestment() {
		return nil
	}

	now := time.Now().UnixMilli()

	// compute latest value and gain
	var totalBalance, totalGain float64
	for _, h := range ac.Holdings {
		lv := h.GetLatestValue()
		gain := h.GetGain()

		if ac.GetCurrency() != h.GetCurrency() {
			erf := repo.NewExchangeRateFilter(
				repo.WithExchangeRateFrom(h.Currency),
				repo.WithExchangeRateTo(ac.Currency),
				repo.WithExchangeRateTimestamp(goutil.Uint64(uint64(now))),
			)
			er, err := uc.exchangeRateRepo.Get(ctx, erf)
			if err != nil {
				return fmt.Errorf("fail to get exchange rate from repo, err: %v", err)
			}

			lv *= er.GetRate()
			gain *= er.GetRate()
		}

		totalBalance += lv
		totalGain += gain
	}
	ac.SetBalance(goutil.Float64(totalBalance))
	ac.SetGain(goutil.Float64(totalGain))

	// compute weighted average percent gain
	var percentGain float64
	if totalBalance > 0 {
		for _, h := range ac.Holdings {
			lv := h.GetLatestValue()

			if ac.GetCurrency() != h.GetCurrency() {
				er, err := uc.exchangeRateRepo.Get(ctx, repo.NewExchangeRateFilter(
					repo.WithExchangeRateFrom(h.Currency),
					repo.WithExchangeRateTo(ac.Currency),
					repo.WithExchangeRateTimestamp(goutil.Uint64(uint64(now))),
				))
				if err != nil {
					return fmt.Errorf("fail to get exchange rate from repo, err: %v", err)
				}
				lv *= er.GetRate()
			}

			weight := lv / totalBalance
			percentGain += weight * h.GetPercentGain()
		}
	}
	ac.SetPercentGain(goutil.Float64(percentGain))

	return nil
}

func (uc *accountUseCase) GetAccountsSummary(ctx context.Context, req *GetAccountsSummaryRequest) (*GetAccountsSummaryResponse, error) {
	user := req.GetUser()

	sf, err := req.ToSnapshotFilter()
	if err != nil {
		return nil, err
	}

	sps, err := uc.snapshotRepo.GetMany(ctx, sf)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get account snapshots from repo, err: %v", err)
		return nil, err
	}

	loc, err := time.LoadLocation(req.AppMeta.GetTimezone())
	if err != nil {
		return nil, entity.ErrInvalidTimezone
	}

	// fetch latest snapshot
	latestGetAccounts, err := uc.GetAccounts(ctx, &GetAccountsRequest{
		UserID: req.User.UserID,
	})
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get latest accounts, err: %v", err)
		return nil, err
	}

	b, err := json.Marshal(latestGetAccounts)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to marshal latest snapshot, err: %v", err)
		return nil, err
	}
	latestSnapshot := string(b)

	now := time.Now().In(loc)

	sps = append(sps, entity.NewSnapshot(
		user.GetUserID(),
		uint32(entity.SnapshotTypeAccount),
		entity.WithSnapshotRecord(goutil.String(latestSnapshot)),
		entity.WithSnapshotTimestamp(goutil.Uint64(uint64(now.UnixMilli()))), // latest
	))

	var (
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

		snapshotsByDate = make(map[string]*GetAccountsResponse)
	)
	for i > 0 {
		date = date.AddDate(yearDiff, monthDiff, 0)
		d := util.FormatDate(date)

		snapshotsByDate[d] = new(GetAccountsResponse) // dummy placeholder

		i--
	}

	var (
		latestDates      = make(map[string]time.Time)
		latestTimestamps = make(map[string]uint64)
	)
	for _, sp := range sps {
		var (
			date time.Time
			t    = time.UnixMilli(int64(sp.GetTimestamp())).In(loc)
		)
		switch req.GetUnit() {
		case uint32(entity.SnapshotUnitMonth):
			date = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()) // start of month
		}

		d := util.FormatDate(date)

		// in case there are multiple snapshots in the same unit, we get the latest snapshot
		if existing, ok := latestDates[d]; !ok || t.After(existing) {
			latestDates[d] = t
			latestTimestamps[d] = sp.GetTimestamp()

			getAccountsResp := new(GetAccountsResponse)
			if err := json.Unmarshal([]byte(sp.GetRecord()), &getAccountsResp); err != nil {
				log.Ctx(ctx).Error().Msgf("fail to unmarshal snapshot, snapshot: %v, err: %v", sp.GetRecord(), err)
				return nil, err
			}

			snapshotsByDate[d] = getAccountsResp
		}
	}

	dates := make([]string, 0)
	for date := range snapshotsByDate {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	resp := new(GetAccountsSummaryResponse)
	for i, date := range dates {
		snapshot, ok := snapshotsByDate[date]
		if !ok {
			snapshot = &GetAccountsResponse{
				NetWorth:   goutil.Float64(0),
				AssetValue: goutil.Float64(0),
				DebtValue:  goutil.Float64(0),
				Currency:   goutil.String(""),
			}
		}

		var (
			netWorth   = snapshot.GetNetWorth()
			assetValue = snapshot.GetAssetValue()
			debtValue  = snapshot.GetDebtValue()
		)
		if snapshot.GetCurrency() != "" && user.Meta.GetCurrency() != snapshot.GetCurrency() {
			erf := repo.NewExchangeRateFilter(
				repo.WithExchangeRateFrom(snapshot.Currency),
				repo.WithExchangeRateTo(user.Meta.Currency),
				repo.WithExchangeRateTimestamp(goutil.Uint64(uint64(latestTimestamps[date]))),
			)
			er, err := uc.exchangeRateRepo.Get(ctx, erf)
			if err != nil {
				return nil, fmt.Errorf("fail to get exchange rate from repo, err: %v", err)
			}

			netWorth *= er.GetRate()
			assetValue *= er.GetRate()
			debtValue *= er.GetRate()
		}

		var (
			percentNetWorthChange, percentAssetChange, percentDebtChange    *float64
			absoluteNetWorthChange, absoluteAssetChange, absoluteDebtChange *float64
		)
		if i == len(dates)-1 {
			var (
				firstNetWorth   = resp.NetWorth[0].GetSum()
				firstAssetValue = resp.AssetValue[0].GetSum()
				firstDebtValue  = resp.DebtValue[0].GetSum()
			)
			absoluteNetWorthChange = goutil.Float64(netWorth - firstNetWorth)
			if *absoluteNetWorthChange == 0 {
				percentNetWorthChange = goutil.Float64(0)
			} else if firstNetWorth != 0 {
				percentNetWorthChange = goutil.Float64(*absoluteNetWorthChange * 100 / math.Abs(firstNetWorth))
			}

			absoluteAssetChange = goutil.Float64(assetValue - firstAssetValue)
			if *absoluteAssetChange == 0 {
				percentAssetChange = goutil.Float64(0)
			} else if firstAssetValue != 0 {
				percentAssetChange = goutil.Float64(*absoluteAssetChange * 100 / math.Abs(firstAssetValue))
			}

			absoluteDebtChange = goutil.Float64(debtValue - firstDebtValue)
			if *absoluteDebtChange == 0 {
				percentDebtChange = goutil.Float64(0)
			} else if firstDebtValue != 0 {
				percentDebtChange = goutil.Float64(*absoluteDebtChange * 100 / math.Abs(firstDebtValue))
			}
		}

		// net worth
		resp.NetWorth = append(resp.NetWorth, common.NewSummary(
			common.WithSummaryDate(goutil.String(date)),
			common.WithSummarySum(goutil.Float64(netWorth)),
			common.WithSummaryCurrency(user.Meta.Currency),
			common.WithSummaryPercentChange(percentNetWorthChange),
			common.WithSummaryAbsoluteChange(absoluteNetWorthChange),
		))

		// asset value
		resp.AssetValue = append(resp.AssetValue, common.NewSummary(
			common.WithSummaryDate(goutil.String(date)),
			common.WithSummarySum(goutil.Float64(assetValue)),
			common.WithSummaryCurrency(user.Meta.Currency),
			common.WithSummaryPercentChange(percentAssetChange),
			common.WithSummaryAbsoluteChange(absoluteAssetChange),
		))

		// debt value
		resp.DebtValue = append(resp.DebtValue, common.NewSummary(
			common.WithSummaryDate(goutil.String(date)),
			common.WithSummarySum(goutil.Float64(debtValue)),
			common.WithSummaryCurrency(user.Meta.Currency),
			common.WithSummaryPercentChange(percentDebtChange),
			common.WithSummaryAbsoluteChange(absoluteDebtChange),
		))
	}

	return resp, nil
}
