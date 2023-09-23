package user

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var GetUserValidator = validator.MustForm(map[string]validator.Validator{})

func (h *userHandler) GetUser(ctx context.Context, req *presenter.GetUserRequest, res *presenter.GetUserResponse) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.userUseCase.GetUser(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get user, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
