package account

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/handler/holding"
	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var CreateAccountValidator = validator.MustForm(map[string]validator.Validator{
	"account_name": &validator.String{
		Optional: false,
	},
	"balance": &validator.String{
		Optional:   true,
		UnsetZero:  true,
		Validators: []validator.StringFunc{entity.CheckMonetaryStr},
	},
	"note": &validator.String{
		Optional: true,
		MaxLen:   uint32(config.MaxAccountNoteLength),
	},
	"account_type": &validator.UInt32{
		Optional:   false,
		Validators: []validator.UInt32Func{entity.CheckChildAccountType},
	},
	"holdings": &validator.Slice{
		Optional:  true,
		MaxLen:    20,
		Validator: holding.NewCreateHoldingValidator(true),
	},
})

func (h *accountHandler) CreateAccount(ctx context.Context, req *presenter.CreateAccountRequest, res *presenter.CreateAccountResponse) error {
	user := entity.GetUserFromCtx(ctx)
	req.Currency = user.Meta.Currency // TODO: Support currency on account creation

	useCaseRes, err := h.accountUseCase.CreateAccount(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to create account, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
