package account

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var GetAccountValidator = validator.MustForm(map[string]validator.Validator{
	"account_id": &validator.String{
		Optional: false,
	},
})

func (h *accountHandler) GetAccount(ctx context.Context, req *presenter.GetAccountRequest, res *presenter.GetAccountResponse) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.accountUseCase.GetAccount(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get account, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
