package account

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var DeleteAccountValidator = validator.MustForm(map[string]validator.Validator{
	"account_id": &validator.String{
		Optional: false,
	},
})

func (h *accountHandler) DeleteAccount(ctx context.Context, req *presenter.DeleteAccountRequest, res *presenter.DeleteAccountResponse) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.accountUseCase.DeleteAccount(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to delete account, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
