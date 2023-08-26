package user

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var SendOTPValidator = validator.MustForm(map[string]validator.Validator{
	"email": &validator.String{
		Optional:   false,
		Validators: []validator.StringFunc{entity.CheckEmail},
	},
})

func (h *userHandler) SendOTP(ctx context.Context, req *presenter.SendOTPRequest, res *presenter.SendOTPResponse) error {
	useCaseRes, err := h.userUseCase.SendOTP(ctx, req.ToUseCaseReq())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to send otp, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
