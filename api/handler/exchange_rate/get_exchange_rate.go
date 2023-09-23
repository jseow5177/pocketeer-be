package exchangerate

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var GetExchangeRateValidator = validator.MustForm(map[string]validator.Validator{
	"from": &validator.String{
		Optional:   false,
		Validators: []validator.StringFunc{entity.CheckCurrency},
	},
	"to": &validator.String{
		Optional:   false,
		Validators: []validator.StringFunc{entity.CheckCurrency},
	},
	"timestamp": &validator.UInt64{
		Optional: false,
	},
})

func (h *exchangeRateHandler) GetExchangeRate(
	ctx context.Context,
	req *presenter.GetExchangeRateRequest,
	res *presenter.GetExchangeRateResponse,
) error {
	useCaseRes, err := h.exchangeRateUseCase.GetExchangeRate(ctx, req.ToUseCaseReq())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get exchange rate, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
