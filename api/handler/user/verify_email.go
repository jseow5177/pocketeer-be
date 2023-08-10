package user

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var VerifyEmailValidator = validator.MustForm(map[string]validator.Validator{
	"email_token": &validator.String{
		Optional: false,
	},
})

func (h *userHandler) VerifyEmail(ctx context.Context, req *presenter.VerifyEmailRequest, res *presenter.VerifyEmailResponse) error {
	useCaseRes, err := h.userUseCase.VerifyEmail(ctx, req.ToUseCaseReq())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to verify email, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
