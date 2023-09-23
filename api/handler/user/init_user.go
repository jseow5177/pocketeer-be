package user

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/handler/account"
	"github.com/jseow5177/pockteer-be/api/handler/category"
	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var InitUserValidator = validator.MustForm(map[string]validator.Validator{
	"currency": &validator.String{
		Optional:   false,
		Validators: []validator.StringFunc{entity.CheckCurrency},
	},
	"categories": &validator.Slice{
		Optional:  true,
		MaxLen:    20,
		Validator: category.CreateCategoryValidator,
	},
	"accounts": &validator.Slice{
		Optional:  true,
		MaxLen:    10,
		Validator: account.CreateAccountValidator,
	},
})

func (h *userHandler) InitUser(ctx context.Context, req *presenter.InitUserRequest, res *presenter.InitUserResponse) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.userUseCase.InitUser(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to init user, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
