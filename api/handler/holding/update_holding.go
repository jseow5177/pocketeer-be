package holding

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var UpdateHoldingValidator = validator.MustForm(map[string]validator.Validator{
	"total_cost": &validator.String{
		Optional:   true,
		UnsetZero:  true,
		Validators: []validator.StringFunc{entity.CheckPositiveMonetaryStr},
	},
	"latest_value": &validator.String{
		Optional:   true,
		UnsetZero:  true,
		Validators: []validator.StringFunc{entity.CheckPositiveMonetaryStr},
	},
})

func (h *holdingHandler) UpdateHolding(ctx context.Context, req *presenter.UpdateHoldingRequest, res *presenter.UpdateHoldingResponse) error {
	userID := util.GetUserIDFromCtx(ctx)

	useCaseRes, err := h.holdingUseCase.UpdateHolding(ctx, req.ToUseCaseReq(userID))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to update holding, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
