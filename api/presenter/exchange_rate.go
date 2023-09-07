package presenter

import (
	exchangerate "github.com/jseow5177/pockteer-be/usecase/exchange_rate"
)

type ExchangeRate struct {
	ExchangeRateID *string  `json:"exchange_rate_id,omitempty"`
	From           *string  `json:"from,omitempty"`
	To             *string  `json:"to,omitempty"`
	Rate           *float64 `json:"rate,omitempty"`
	Timestamp      *uint64  `json:"timestamp,omitempty"`
	CreateTime     *uint64  `json:"create_time,omitempty"`
}

func (er *ExchangeRate) GetExchangeRateID() string {
	if er != nil && er.ExchangeRateID != nil {
		return *er.ExchangeRateID
	}
	return ""
}

func (er *ExchangeRate) GetFrom() string {
	if er != nil && er.From != nil {
		return *er.From
	}
	return ""
}

func (er *ExchangeRate) GetTo() string {
	if er != nil && er.To != nil {
		return *er.To
	}
	return ""
}

func (er *ExchangeRate) GetRate() float64 {
	if er != nil && er.Rate != nil {
		return *er.Rate
	}
	return 0
}

func (er *ExchangeRate) GetTimestamp() uint64 {
	if er != nil && er.Timestamp != nil {
		return *er.Timestamp
	}
	return 0
}

func (er *ExchangeRate) GetCreateTime() uint64 {
	if er != nil && er.CreateTime != nil {
		return *er.CreateTime
	}
	return 0
}

type CreateExchangeRatesRequest struct {
	Date *string  `json:"date,omitempty"`
	From *string  `json:"from,omitempty"`
	To   []string `json:"to,omitempty"`
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

func (m *CreateExchangeRatesRequest) ToUseCaseReq() *exchangerate.CreateExchangeRatesRequest {
	return &exchangerate.CreateExchangeRatesRequest{
		Date: m.Date,
		From: m.From,
		To:   m.To,
	}
}

type CreateExchangeRatesResponse struct {
	ExchangeRates []*ExchangeRate `json:"exchange_rates,omitempty"`
}

func (m *CreateExchangeRatesResponse) GetExchangeRates() []*ExchangeRate {
	if m != nil && m.ExchangeRates != nil {
		return m.ExchangeRates
	}
	return nil
}

func (m *CreateExchangeRatesResponse) Set(useCaseRes *exchangerate.CreateExchangeRatesResponse) {
	m.ExchangeRates = toExchangeRates(useCaseRes.ExchangeRates)
}
