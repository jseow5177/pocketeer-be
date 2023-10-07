package account

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
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

		if err := uc.holdingRepo.DeleteMany(txCtx, req.ToHoldingFilter()); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to mark account holdings as deleted, err: %v", err)
			return err
		}

		if err := uc.lotRepo.DeleteMany(txCtx, req.ToLotFilter(holdingIDs)); err != nil {
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
		account.GetAccountID(),
		"",
		entity.WithTransactionAmount(goutil.Float64(amount)),
		entity.WithTransactionCurrency(account.Currency),
		entity.WithTransactionType(goutil.Uint32(tt)),
		entity.WithTransactionNote(goutil.String(note)),
		entity.WithTransactionTime(goutil.Uint64(uint64(time.Now().UnixMilli()))),
	)
}

func (uc *accountUseCase) getAccountHoldingsAndLots(ctx context.Context, ac *entity.Account) error {
	hs, err := uc.holdingRepo.GetMany(ctx, repo.NewHoldingFilter(
		ac.GetUserID(),
		repo.WithHoldingAccountID(ac.AccountID),
	))
	if err != nil {
		return fmt.Errorf("fail to get holdings from repo, err: %v", err)
	}

	for _, h := range hs {
		if h.IsDefault() {
			q, err := uc.quoteRepo.Get(ctx, &repo.QuoteFilter{
				Symbol: h.Symbol,
			})
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
