package account

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var GetAccountsValidator = validator.MustForm(map[string]validator.Validator{
	"account_type": &validator.UInt32{
		Optional:   true,
		UnsetZero:  true,
		Validators: []validator.UInt32Func{entity.CheckAccountType},
	},
})

func (h *accountHandler) GetAccounts(ctx context.Context, req *presenter.GetAccountsRequest, res *presenter.GetAccountsResponse) error {
	userID := util.GetUserIDFromCtx(ctx)

	useCaseRes, err := h.accountUseCase.GetAccounts(ctx, req.ToUseCaseReq(userID))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get accounts from repo, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
