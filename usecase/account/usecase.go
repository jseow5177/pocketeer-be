package account

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
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

func (uc *accountUseCase) UpdateAccount(ctx context.Context, req *UpdateAccountRequest) (*UpdateAccountResponse, error) {
	acRes, err := uc.GetAccount(ctx, req.ToGetAccountRequest())
	if err != nil {
		return nil, err
	}
	ac := acRes.Account

	nac := uc.getAccountUpdates(ac, req.ToAccountEntity())
	if nac == nil {
		// no updates
		log.Ctx(ctx).Info().Msg("acount has no updates")
		return &UpdateAccountResponse{
			Account: ac,
		}, nil
	}

	if err = uc.accountRepo.Update(ctx, req.ToAccountFilter(), nac); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save account updates to repo, err: %v", err)
		return nil, err
	}

	// merge
	goutil.MergeWithPtrFields(ac, nac)

	return &UpdateAccountResponse{
		Account: ac,
	}, nil
}

func (uc *accountUseCase) getAccountUpdates(old, changes *entity.Account) *entity.Account {
	var hasUpdates bool

	nac := new(entity.Account)

	if changes.AccountName != nil && changes.GetAccountName() != old.GetAccountName() {
		hasUpdates = true
		nac.AccountName = changes.AccountName
	}

	if changes.Balance != nil && changes.GetBalance() != old.GetBalance() {
		hasUpdates = true
		nac.Balance = changes.Balance
	}

	if changes.Note != nil && changes.GetNote() != old.GetNote() {
		hasUpdates = true
		nac.Note = changes.Note
	}

	if !hasUpdates {
		return nil
	}

	return nac
}
