package user

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var SignUpValidator = validator.MustForm(map[string]validator.Validator{
	"username": &validator.String{
		Optional: false,
		MaxLen:   config.UsernameMaxLength,
	},
	"password": &validator.String{
		Optional: false,
		MinLen:   config.PasswordMinLength,
	},
})

func (h *userHandler) SignUp(ctx context.Context, req *presenter.SignUpRequest, res *presenter.SignUpResponse) error {
	useCaseRes, err := h.userUseCase.SignUp(ctx, req.ToUseCaseReq())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to sign up user, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
