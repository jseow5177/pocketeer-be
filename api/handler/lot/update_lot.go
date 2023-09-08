package lot

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var UpdateLotValidator = validator.MustForm(map[string]validator.Validator{
	"lot_id": &validator.String{
		Optional: false,
	},
	"shares": &validator.String{
		Optional:   true,
		UnsetZero:  true,
		Validators: []validator.StringFunc{entity.CheckPositiveMonetaryStr},
	},
	"cost_per_share": &validator.String{
		Optional:   true,
		UnsetZero:  true,
		Validators: []validator.StringFunc{entity.CheckPositiveMonetaryStr},
	},
	"trade_date": &validator.UInt64{
		Optional:  true,
		UnsetZero: true,
	},
})

func (h *lotHandler) UpdateLot(ctx context.Context, req *presenter.UpdateLotRequest, res *presenter.UpdateLotResponse) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.lotUseCase.UpdateLot(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to update lot, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
