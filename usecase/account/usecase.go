package account

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/rs/zerolog/log"
)

type accountUseCase struct {
	txMgr           repo.TxMgr
	accountRepo     repo.AccountRepo
	transactionRepo repo.TransactionRepo
}

func NewAccountUseCase(txMgr repo.TxMgr, accountRepo repo.AccountRepo, transactionRepo repo.TransactionRepo) UseCase {
	return &accountUseCase{
		txMgr,
		accountRepo,
		transactionRepo,
	}
}

func (uc *accountUseCase) GetAccount(ctx context.Context, req *GetAccountRequest) (*GetAccountResponse, error) {
	ac, err := uc.accountRepo.Get(ctx, req.ToAccountFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get account from repo, err: %v", err)
		return nil, err
	}

	return &GetAccountResponse{
		Account: ac,
	}, nil
}

func (uc *accountUseCase) CreateAccount(ctx context.Context, req *CreateAccountRequest) (*CreateAccountResponse, error) {
	ac := req.ToAccountEntity()

	_, err := uc.accountRepo.Create(ctx, ac)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save new account to repo, err: %v", err)
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

	nac := ac.GetUpdates(req.ToAccountUpdate(), true)
	if nac == nil {
		log.Ctx(ctx).Info().Msg("acount has no updates")
		return &UpdateAccountResponse{
			Account: ac,
		}, nil
	}

	if err = uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		if err = uc.accountRepo.Update(ctx, req.ToAccountFilter(), nac); err != nil {
			log.Ctx(ctx).Error().Msgf("fail to save account updates to repo, err: %v", err)
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &UpdateAccountResponse{
		Account: ac,
	}, nil
}
