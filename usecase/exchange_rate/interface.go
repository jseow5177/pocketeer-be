package exchangerate

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/api"
	"github.com/jseow5177/pockteer-be/entity"
)

type UseCase interface {
	CreateExchangeRate(ctx context.Context, req *CreateExchangeRatesRequest) (*CreateExchangeRatesResponse, error)
}

type CreateExchangeRatesRequest struct {
	Date *string
	From *string
	To   []string
}

func (m *CreateExchangeRatesRequest) GetFrom() string {
	if m != nil && m.From != nil {
		return *m.From
	}
	return ""
}

func (m *CreateExchangeRatesRequest) GetDate() string {
	if m != nil && m.Date != nil {
		return *m.Date
	}
	return ""
}

func (m *CreateExchangeRatesRequest) GetTo() []string {
	if m != nil && m.To != nil {
		return m.To
	}
	return nil
}

func (m *CreateExchangeRatesRequest) ToExchangeRateFilter() *api.ExchangeRateFilter {
	return &api.ExchangeRateFilter{
		Base:    m.From,
		Symbols: m.To,
		Date:    m.Date,
	}
}

type CreateExchangeRatesResponse struct {
	ExchangeRates []*entity.ExchangeRate
}

func (m *CreateExchangeRatesResponse) GetExchangeRates() []*entity.ExchangeRate {
	if m != nil && m.ExchangeRates != nil {
		return m.ExchangeRates
	}
	return nil
}
