package user

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
)

var GetUserValidator = validator.MustForm(map[string]validator.Validator{
	"user_id": &validator.String{
		Optional: false,
	},
})

func (h *userHandler) GetUser(ctx context.Context, req *presenter.GetUserRequest, res *presenter.GetUserResponse) error {
	return nil
}
