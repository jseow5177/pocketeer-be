package account

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"

	auc "github.com/jseow5177/pockteer-be/usecase/account"
)

var UpdateAccountValidator = validator.MustForm(map[string]validator.Validator{
	"account_name": &validator.String{
		Optional:  true,
		UnsetZero: true,
	},
	"balance": &validator.String{
		Optional:   true,
		UnsetZero:  true,
		Validators: []validator.StringFunc{entity.CheckMonetaryStr},
	},
	"note": &validator.String{
		Optional:  true,
		UnsetZero: true,
		MaxLen:    uint32(config.MaxAccountNoteLength),
	},
	"update_mode": &validator.UInt32{
		Optional:   true,
		Validators: []validator.UInt32Func{auc.CheckUpdateMode},
	},
})

func (h *accountHandler) UpdateAccount(ctx context.Context, req *presenter.UpdateAccountRequest, res *presenter.UpdateAccountResponse) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.accountUseCase.UpdateAccount(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to update account, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
