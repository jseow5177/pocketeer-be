package user

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var UpdateUserMetaValidator = validator.MustForm(map[string]validator.Validator{
	"init_stage": &validator.UInt32{
		Optional:  true,
		UnsetZero: true,
		Min:       goutil.Uint32(uint32(entity.InitStageOne)),
		Max:       goutil.Uint32(uint32(entity.InitStageFour)),
	},
})

func (h *userHandler) UpdateUserMeta(ctx context.Context, req *presenter.UpdateUserMetaRequest, res *presenter.UpdateUserMetaResponse) error {
	userID := util.GetUserIDFromCtx(ctx)

	useCaseRes, err := h.userUseCase.UpdateUserMeta(ctx, req.ToUseCaseReq(userID))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to update user meta, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
