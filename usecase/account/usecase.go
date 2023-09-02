package account

import (
	"context"
	"fmt"
	"time"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/rs/zerolog/log"
)

type accountUseCase struct {
	txMgr           repo.TxMgr
	accountRepo     repo.AccountRepo
	transactionRepo repo.TransactionRepo
	holdingRepo     repo.HoldingRepo
	lotRepo         repo.LotRepo
	quoteRepo       repo.QuoteRepo
	securityRepo    repo.SecurityRepo
}

func NewAccountUseCase(
	txMgr repo.TxMgr,
	accountRepo repo.AccountRepo,
	transactionRepo repo.TransactionRepo,
	holdingRepo repo.HoldingRepo,
	lotRepo repo.LotRepo,
	quoteRepo repo.QuoteRepo,
	securityRepo repo.SecurityRepo,
) UseCase {
	return &accountUseCase{
		txMgr,
		accountRepo,
		transactionRepo,
		holdingRepo,
		lotRepo,
		quoteRepo,
		securityRepo,
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

	if err := goutil.ParallelizeWork(ctx, len(acs), 5, func(ctx context.Context, workNum int) error {
		ac := acs[workNum]

		if !ac.IsInvestment() {
			return nil
		}

		if err := uc.getAccountHoldingsAndLots(ctx, ac); err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get account holdings and lots, err: %v", err)
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &GetAccountsResponse{
		Accounts: acs,
	}, nil
}

func (uc *accountUseCase) CreateAccount(ctx context.Context, req *CreateAccountRequest) (*CreateAccountResponse, error) {
	ac, err := req.ToAccountEntity()
	if err != nil {
		return nil, err
	}

	if err := uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		if _, err := uc.accountRepo.Create(txCtx, ac); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to save new account to repo, err: %v", err)
			return err
		}

		hs := ac.Holdings
		if len(hs) == 0 {
			return nil
		}

		for _, h := range hs {
			h.SetAccountID(ac.AccountID)
		}

		for i, h := range hs {
			if h.IsDefault() {
				if _, err = uc.securityRepo.Get(txCtx, req.Holdings[i].ToSecurityFilter()); err != nil {
					log.Ctx(txCtx).Error().Msgf("fail to get security from repo, err: %v", err)
					return fmt.Errorf("symbol %v, err: %v", h.GetSymbol(), err)
				}

				q, err := uc.quoteRepo.Get(txCtx, req.Holdings[i].ToQuoteFilter())
				if err != nil {
					log.Ctx(txCtx).Error().Msgf("fail to get quote from repo, err: %v", err)
					return err
				}
				h.SetQuote(q)
			}
		}

		_, err = uc.holdingRepo.CreateMany(txCtx, hs)
		if err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to save new holdings to repo, err: %v", err)
			return err
		}

		if err := goutil.ParallelizeWork(txCtx, len(hs), 5, func(ctx context.Context, workNum int) error {
			h := hs[workNum]

			if len(h.Lots) == 0 {
				return nil
			}

			ls := h.Lots
			for _, l := range ls {
				l.SetHoldingID(h.HoldingID)
			}

			_, err := uc.lotRepo.CreateMany(ctx, ls)
			if err != nil {
				log.Ctx(ctx).Error().Msgf("fail to save new lots to repo, err: %v", err)
				return err
			}

			h.SetLots(ls)
			h.ComputeSharesCostAndValue()

			return nil
		}); err != nil {
			return err
		}

		ac.SetHoldings(hs)
		ac.ComputeCostAndBalance()

		return nil
	}); err != nil {
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

	acu, err := ac.Update(req.ToAccountUpdate())
	if err != nil {
		return nil, err
	}

	if acu == nil {
		log.Ctx(ctx).Info().Msg("acount has no updates")
		return &UpdateAccountResponse{
			Account: ac,
		}, nil
	}

	if err = uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		if err = uc.accountRepo.Update(txCtx, req.ToAccountFilter(), acu); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to save account updates to repo, err: %v", err)
			return err
		}

		if acu.Balance != nil && req.NeedOffsetTransaction() {
			balanceChange := acu.GetBalance() - oldBalance
			t := uc.newUnrecordedTransaction(balanceChange, ac.GetUserID(), ac.GetAccountID())

			if _, err = uc.transactionRepo.Create(txCtx, t); err != nil {
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

func (uc *accountUseCase) newUnrecordedTransaction(amount float64, userID, accountID string) *entity.Transaction {
	tt := uint32(entity.GetTransactionTypeByAmount(amount))

	note := fmt.Sprintf("Unrecorded %s", entity.TransactionTypes[tt])

	t := entity.NewTransaction(
		userID,
		accountID,
		"",
		entity.WithTransactionAmount(goutil.Float64(amount)),
		entity.WithTransactionType(goutil.Uint32(tt)),
		entity.WithTransactionNote(goutil.String(note)),
		entity.WithTransactionTime(goutil.Uint64(uint64(time.Now().UnixMilli()))),
	)

	return t
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
		if !h.IsDefault() {
			continue
		}

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

		h.ComputeSharesCostAndValue()
	}

	ac.SetHoldings(hs)
	ac.ComputeCostAndBalance()

	return nil
}
