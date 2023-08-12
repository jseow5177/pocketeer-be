package user

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var InitUserValidator = validator.MustForm(map[string]validator.Validator{})

func (h *userHandler) InitUser(ctx context.Context, req *presenter.InitUserRequest, res *presenter.InitUserResponse) error {
	userID := util.GetUserIDFromCtx(ctx)

	useCaseRes, err := h.userUseCase.InitUser(ctx, req.ToUseCaseReq(userID))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to init user, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
