package user

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/pkg/validator"
)

var LogInValidator = validator.MustForm(map[string]validator.Validator{
	"username": &validator.String{
		Optional: false,
		MaxLen:   config.UsernameMaxLength,
	},
	"password": &validator.String{
		Optional: false,
		MinLen:   config.PasswordMinLength,
	},
})

func (h *userHandler) LogIn(ctx context.Context, req *presenter.LogInRequest, res *presenter.LogInResponse) error {
	return nil
}
