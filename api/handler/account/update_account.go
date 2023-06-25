package account

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var UpdateAccountValidator = validator.MustForm(map[string]validator.Validator{
	"account_name": &validator.String{
		Optional: true,
	},
	"balance": &validator.String{
		Optional:   true,
		Validators: []validator.StringFunc{entity.CheckMonetaryStr},
	},
	"note": &validator.String{
		Optional: true,
		MaxLen:   uint32(config.MaxAccountNoteLength),
	},
})

func (h *accountHandler) UpdateAccount(ctx context.Context, req *presenter.UpdateAccountRequest, res *presenter.UpdateAccountResponse) error {
	userID := util.GetUserIDFromCtx(ctx)

	useCaseRes, err := h.accountUseCase.UpdateAccount(ctx, req.ToUseCaseReq(userID))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to update account, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
