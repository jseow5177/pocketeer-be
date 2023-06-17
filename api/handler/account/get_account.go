package account

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var GetAccountValidator = validator.MustForm(map[string]validator.Validator{
	"account_id": &validator.String{
		Optional: false,
	},
})

func (h *accountHandler) GetAccount(ctx context.Context, req *presenter.GetAccountRequest, res *presenter.GetAccountResponse) error {
	userID := util.GetUserIDFromCtx(ctx)

	useCaseRes, err := h.accountUseCase.GetAccount(ctx, req.ToUseCaseReq(userID))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get ACCOUNT, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
