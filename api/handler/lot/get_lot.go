package lot

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var GetLotValidator = validator.MustForm(map[string]validator.Validator{
	"lot_id": &validator.String{
		Optional: false,
	},
})

func (h *lotHandler) GetLot(ctx context.Context, req *presenter.GetLotRequest, res *presenter.GetLotResponse) error {
	userID := util.GetUserIDFromCtx(ctx)

	useCaseRes, err := h.lotUseCase.GetLot(ctx, req.ToUseCaseReq(userID))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get lot, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
