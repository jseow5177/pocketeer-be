package user

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var UpdateUserMetaValidator = validator.MustForm(map[string]validator.Validator{
	"currency": &validator.String{
		Optional:   true,
		UnsetZero:  true,
		Validators: []validator.StringFunc{entity.CheckCurrency},
	},
})

func (h *userHandler) UpdateUserMeta(ctx context.Context, req *presenter.UpdateUserMetaRequest, res *presenter.UpdateUserMetaResponse) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.userUseCase.UpdateUserMeta(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to update user meta, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
