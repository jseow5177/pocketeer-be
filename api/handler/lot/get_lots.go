package lot

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var GetLotsValidator = validator.MustForm(map[string]validator.Validator{
	"holding_id": &validator.String{
		Optional: false,
	},
})

func (h *lotHandler) GetLots(ctx context.Context, req *presenter.GetLotsRequest, res *presenter.GetLotsResponse) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.lotUseCase.GetLots(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get lots from repo, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
