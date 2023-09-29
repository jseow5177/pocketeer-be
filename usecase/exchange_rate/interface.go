package exchangerate

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
)

type UseCase interface {
	GetCurrencies(ctx context.Context, req *GetCurrenciesRequest) (*GetCurrenciesResponse, error)
	GetExchangeRate(ctx context.Context, req *GetExchangeRateRequest) (*GetExchangeRateResponse, error)
}

type GetCurrenciesRequest struct {
	User *entity.User
}

func (m *GetCurrenciesRequest) GetUser() *entity.User {
	if m != nil && m.User != nil {
		return m.User
	}
	return nil
}

type GetCurrenciesResponse struct {
	Currencies []string
}

func (m *GetCurrenciesResponse) GetCurrencies() []string {
	if m != nil && m.Currencies != nil {
		return m.Currencies
	}
	return nil
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

func (m *GetExchangeRateRequest) ToGetExchangeRateFilter() *repo.ExchangeRateFilter {
	return repo.NewExchangeRateFilter(
		repo.WithExchangeRateFrom(m.From),
		repo.WithExchangeRateTo(m.To),
		repo.WithExchangeRateTimestamp(m.Timestamp),
	)
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
