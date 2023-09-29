package exchangerate

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var GetCurrenciesValidator = validator.MustForm(map[string]validator.Validator{})

func (h *exchangeRateHandler) GetCurrencies(
	ctx context.Context,
	req *presenter.GetCurrenciesRequest,
	res *presenter.GetCurrenciesResponse,
) error {
	u := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.exchangeRateUseCase.GetCurrencies(ctx, req.ToUseCaseReq(u))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get currencies, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
