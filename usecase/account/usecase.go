package account

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/holding"
	"github.com/rs/zerolog/log"
)

type accountUseCase struct {
	txMgr           repo.TxMgr
	accountRepo     repo.AccountRepo
	transactionRepo repo.TransactionRepo
	holdingUseCase  holding.UseCase
}

func NewAccountUseCase(
	txMgr repo.TxMgr,
	accountRepo repo.AccountRepo,
	transactionRepo repo.TransactionRepo,
	holdingUseCase holding.UseCase,
) UseCase {
	return &accountUseCase{
		txMgr,
		accountRepo,
		transactionRepo,
		holdingUseCase,
	}
}

func (uc *accountUseCase) GetAccount(ctx context.Context, req *GetAccountRequest) (*GetAccountResponse, error) {
	ac, err := uc.accountRepo.Get(ctx, req.ToAccountFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get account from repo, err: %v", err)
		return nil, err
	}

	if ac.IsInvestment() {
		if err = uc.calcInvestmentAccountValue(ctx, ac); err != nil {
			log.Ctx(ctx).Error().Msgf("fail to compute account value, err: %v", err)
			return nil, err
		}
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
		if ac.IsInvestment() {
			return uc.calcInvestmentAccountValue(ctx, acs[workNum])
		}
		return nil
	}); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to compute accounts value, err: %v", err)
		return nil, err
	}

	return &GetAccountsResponse{
		Accounts: acs,
	}, nil
}

func (uc *accountUseCase) CreateAccounts(ctx context.Context, req *CreateAccountsRequest) (*CreateAccountsResponse, error) {
	var (
		mu      sync.Mutex
		acs     = make([]*entity.Account, 0)
		errChan = make(chan int, len(req.Accounts))
	)

	if err := uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		return goutil.ParallelizeWork(ctx, len(req.Accounts), 5, func(ctx context.Context, workNum int) error {
			acRes, err := uc.CreateAccount(ctx, req.Accounts[workNum])
			if err != nil {
				errChan <- workNum
				return err
			}
			mu.Lock()
			acs = append(acs, acRes.Account)
			mu.Unlock()
			return nil
		})
	}); err != nil {
		return nil, fmt.Errorf("account idx %v, err: %v", <-errChan, err)
	}

	return &CreateAccountsResponse{
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

		if len(req.Holdings) > 0 {
			holdingsRes, err := uc.holdingUseCase.CreateHoldings(txCtx, req.ToCreateHoldingsRequest(ac.GetAccountID()))
			if err != nil {
				return err
			}
			ac.SetHoldings(holdingsRes.Holdings)
		}

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

	acu, hasUpdate, err := ac.Update(req.ToAccountUpdate())
	if err != nil {
		return nil, err
	}

	if !hasUpdate {
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

func (uc *accountUseCase) calcInvestmentAccountValue(ctx context.Context, ac *entity.Account) error {
	res, err := uc.holdingUseCase.GetHoldings(ctx, &holding.GetHoldingsRequest{
		UserID:    ac.UserID,
		AccountID: ac.AccountID,
	})
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get account holdings, err: %v", err)
		return err
	}
	ac.SetHoldings(res.Holdings)

	var (
		totalCost   float64
		latestValue float64
	)
	for _, h := range res.Holdings {
		totalCost += h.GetTotalCost()
		latestValue += h.GetLatestValue()
	}
	ac.SetBalance(goutil.Float64(latestValue)) // store latest value into balance field
	ac.SetTotalCost(goutil.Float64(totalCost))

	return nil
}
