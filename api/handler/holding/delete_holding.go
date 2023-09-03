package holding

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var DeleteHoldingValidator = validator.MustForm(map[string]validator.Validator{
	"holding_id": &validator.String{
		Optional: false,
	},
})

func (h *holdingHandler) DeleteHolding(ctx context.Context, req *presenter.DeleteHoldingRequest, res *presenter.DeleteHoldingResponse) error {
	userID := util.GetUserIDFromCtx(ctx)

	useCaseRes, err := h.holdingUseCase.DeleteHolding(ctx, req.ToUseCaseReq(userID))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to delete holding, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
