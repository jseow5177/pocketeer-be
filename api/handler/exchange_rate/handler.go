package exchangerate

import (
	exchangerate "github.com/jseow5177/pockteer-be/usecase/exchange_rate"
)

type exchangeRateHandler struct {
	exchangeRateUseCase exchangerate.UseCase
}

func NewExchangeRateHandler(exchangeRateUseCase exchangerate.UseCase) *exchangeRateHandler {
	return &exchangeRateHandler{
		exchangeRateUseCase,
	}
}
