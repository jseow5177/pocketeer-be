package lot

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var CreateLotValidator = validator.MustForm(map[string]validator.Validator{
	"holding_id": &validator.String{
		Optional: false,
	},
	"shares": &validator.String{
		Optional:   false,
		Validators: []validator.StringFunc{entity.CheckMonetaryStr},
	},
	"cost_per_share": &validator.String{
		Optional:   false,
		Validators: []validator.StringFunc{entity.CheckMonetaryStr},
	},
	"trade_date": &validator.UInt64{
		Optional: false,
	},
})

func (h *lotHandler) CreateLot(ctx context.Context, req *presenter.CreateLotRequest, res *presenter.CreateLotResponse) error {
	userID := util.GetUserIDFromCtx(ctx)

	useCaseRes, err := h.lotUseCase.CreateLot(ctx, req.ToUseCaseReq(userID))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to create lot, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
