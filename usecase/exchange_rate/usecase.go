package exchangerate

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/api"
	"github.com/jseow5177/pockteer-be/dep/repo"
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

func (uc *exchangeRateUseCase) CreateExchangeRates(ctx context.Context, req *CreateExchangeRatesRequest) (*CreateExchangeRatesResponse, error) {
	ers, err := uc.exchangeRateAPI.GetExchangeRates(ctx, req.ToExchangeRateFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get exchange rates, err: %v", err)
		return nil, err
	}

	if _, err := uc.exchangeRateRepo.CreateMany(ctx, ers); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to create exchange rates in repo, err: %v", err)
		return nil, err
	}

	return &CreateExchangeRatesResponse{
		ExchangeRates: ers,
	}, nil
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
