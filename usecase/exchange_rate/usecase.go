package exchangerate

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/api"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/rs/zerolog/log"
)

type exchangeRateUseCase struct {
	exchangeRateAPI  api.ExchangeRateAPI
	exchangeRateRepo repo.ExchangeRateRepo
}

func NewExchangeRateUseCase(
	exchangeRateAPI api.ExchangeRateAPI,
	exchangeRateRepo repo.ExchangeRateRepo,
) UseCase {
	return &exchangeRateUseCase{
		exchangeRateAPI,
		exchangeRateRepo,
	}
}

func (uc *exchangeRateUseCase) GetExchangeRate(ctx context.Context, req *GetExchangeRateRequest) (*GetExchangeRateResponse, error) {
	er, err := uc.exchangeRateRepo.Get(ctx, req.ToGetExchangeRateFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get exchange rate from repo, err: %v", err)
		return nil, err
	}

	return &GetExchangeRateResponse{
		ExchangeRate: er,
	}, nil
}

func (uc *exchangeRateUseCase) GetCurrencies(ctx context.Context, req *GetCurrenciesRequest) (*GetCurrenciesResponse, error) {
	u := req.GetUser()

	currencies := make([]string, 0)
	currencies = append(currencies, u.Meta.GetCurrency())

	for _, currency := range entity.Currencies {
		if currency != u.Meta.GetCurrency() {
			currencies = append(currencies, currency)
		}
	}

	return &GetCurrenciesResponse{
		Currencies: currencies,
	}, nil
}
