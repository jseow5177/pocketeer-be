package exchangerate

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var CreateExchangeRatesValidator = validator.MustForm(map[string]validator.Validator{
	"from": &validator.String{
		Optional:   false,
		Validators: []validator.StringFunc{entity.CheckCurrency},
	},
	"to": &validator.Slice{
		Optional: false,
		MinLen:   1,
		Validator: &validator.String{
			Validators: []validator.StringFunc{entity.CheckCurrency},
		},
	},
	"date": &validator.String{
		Optional:   true,
		UnsetZero:  true,
		Validators: []validator.StringFunc{entity.CheckDateStr},
	},
})

func (h *exchangeRateHandler) CreateExchangeRates(
	ctx context.Context,
	req *presenter.CreateExchangeRatesRequest,
	res *presenter.CreateExchangeRatesResponse,
) error {
	useCaseRes, err := h.exchangeRateUseCase.CreateExchangeRate(ctx, req.ToUseCaseReq())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to create exchange rates, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
