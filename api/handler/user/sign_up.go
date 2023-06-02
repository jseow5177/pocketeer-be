package user

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/pkg/validator"
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
	return nil
}
