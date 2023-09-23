package holding

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var GetHoldingValidator = validator.MustForm(map[string]validator.Validator{
	"holding_id": &validator.String{
		Optional: false,
	},
})

func (h *holdingHandler) GetHolding(ctx context.Context, req *presenter.GetHoldingRequest, res *presenter.GetHoldingResponse) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.holdingUseCase.GetHolding(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get holding, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
