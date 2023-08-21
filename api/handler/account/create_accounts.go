package account

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var CreateAccountsValidator = validator.MustForm(map[string]validator.Validator{
	"accounts": &validator.Slice{
		Optional:  false,
		MaxLen:    10,
		Validator: CreateAccountValidator,
	},
})

func (h *accountHandler) CreateAccounts(ctx context.Context, req *presenter.CreateAccountsRequest, res *presenter.CreateAccountsResponse) error {
	userID := util.GetUserIDFromCtx(ctx)

	useCaseRes, err := h.accountUseCase.CreateAccounts(ctx, req.ToUseCaseReq(userID))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to create accounts, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
