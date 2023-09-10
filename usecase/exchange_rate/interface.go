package exchangerate

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/api"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
)

type UseCase interface {
	CreateExchangeRates(ctx context.Context, req *CreateExchangeRatesRequest) (*CreateExchangeRatesResponse, error)
	GetExchangeRate(ctx context.Context, req *GetExchangeRateRequest) (*GetExchangeRateResponse, error)
}

type GetExchangeRateRequest struct {
	Timestamp *uint64
	From      *string
	To        *string
}

func (m *GetExchangeRateRequest) GetFrom() string {
	if m != nil && m.From != nil {
		return *m.From
	}
	return ""
}

func (m *GetExchangeRateRequest) GetTimestamp() uint64 {
	if m != nil && m.Timestamp != nil {
		return *m.Timestamp
	}
	return 0
}

func (m *GetExchangeRateRequest) GetTo() string {
	if m != nil && m.To != nil {
		return *m.To
	}
	return ""
}

func (m *GetExchangeRateRequest) ToGetExchangeRateFilter() *repo.GetExchangeRateFilter {
	return &repo.GetExchangeRateFilter{
		From:      m.From,
		To:        m.To,
		Timestamp: m.Timestamp,
	}
}

type GetExchangeRateResponse struct {
	ExchangeRate *entity.ExchangeRate
}

func (m *GetExchangeRateResponse) GetExchangeRate() *entity.ExchangeRate {
	if m != nil && m.ExchangeRate != nil {
		return m.ExchangeRate
	}
	return nil
}

type CreateExchangeRatesRequest struct {
	Timestamp *uint64
	From      *string
	To        []string
}

func (m *CreateExchangeRatesRequest) GetFrom() string {
	if m != nil && m.From != nil {
		return *m.From
	}
	return ""
}

func (m *CreateExchangeRatesRequest) GetTimestamp() uint64 {
	if m != nil && m.Timestamp != nil {
		return *m.Timestamp
	}
	return 0
}

func (m *CreateExchangeRatesRequest) GetTo() []string {
	if m != nil && m.To != nil {
		return m.To
	}
	return nil
}

func (m *CreateExchangeRatesRequest) ToExchangeRateFilter() *api.ExchangeRateFilter {
	return api.NewExchangeRateFilter(m.Timestamp, m.GetFrom(), m.To...)
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
