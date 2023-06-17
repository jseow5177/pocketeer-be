package account

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/rs/zerolog/log"
)

type accountUseCase struct {
	accountRepo repo.AccountRepo
}

func NewAccountUseCase(accountRepo repo.AccountRepo) UseCase {
	return &accountUseCase{
		accountRepo,
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

	id, err := uc.accountRepo.Create(ctx, ac)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save new account to repo, err: %v", err)
		return nil, err
	}

	ac.AccountID = goutil.String(id)

	return &CreateAccountResponse{
		Account: ac,
	}, nil
}
