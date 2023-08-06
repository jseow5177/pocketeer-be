package lot

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var DeleteLotValidator = validator.MustForm(map[string]validator.Validator{
	"lot_id": &validator.String{
		Optional: false,
	},
})

func (h *lotHandler) DeleteLot(ctx context.Context, req *presenter.DeleteLotRequest, res *presenter.DeleteLotResponse) error {
	userID := util.GetUserIDFromCtx(ctx)

	useCaseRes, err := h.lotUseCase.DeleteLot(ctx, req.ToUseCaseReq(userID))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to delete lot, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
