package entity

import (
	"time"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/util"
)

type Currency string

const (
	CurrencySGD Currency = "SGD"
	CurrencyUSD Currency = "USD"
	CurrencyRM  Currency = "MYR"
)

var Currencies = map[string]string{
	string(CurrencySGD): "SGD",
	string(CurrencyUSD): "USD",
	string(CurrencyRM):  "MYR",
}

type ExchangeRate struct {
	ExchangeRateID *string
	From           *string
	To             *string
	Rate           *float64
	Timestamp      *uint64
	CreateTime     *uint64
}

type ExchangeRateOption = func(er *ExchangeRate)

func WithExchangeRateID(exchangeRateID *string) ExchangeRateOption {
	return func(er *ExchangeRate) {
		er.SetExchangeRateID(exchangeRateID)
	}
}

func WithExchangeRateCreateTime(createTime *uint64) ExchangeRateOption {
	return func(er *ExchangeRate) {
		er.SetCreateTime(createTime)
	}
}

func NewExchangeRate(from, to string, rate float64, timestamp uint64, opts ...ExchangeRateOption) *ExchangeRate {
	now := uint64(time.Now().UnixMilli())
	er := &ExchangeRate{
		From:       goutil.String(from),
		To:         goutil.String(to),
		Rate:       goutil.Float64(util.RoundFloatToPreciseDP(rate)),
		Timestamp:  goutil.Uint64(timestamp),
		CreateTime: goutil.Uint64(now),
	}
	for _, opt := range opts {
		opt(er)
	}

	return er
}

func (er *ExchangeRate) GetExchangeRateID() string {
	if er != nil && er.ExchangeRateID != nil {
		return *er.ExchangeRateID
	}
	return ""
}

func (er *ExchangeRate) SetExchangeRateID(exchangeRateID *string) {
	er.ExchangeRateID = exchangeRateID
}

func (er *ExchangeRate) GetFrom() string {
	if er != nil && er.From != nil {
		return *er.From
	}
	return ""
}

func (er *ExchangeRate) SetFrom(from *string) {
	er.From = from
}

func (er *ExchangeRate) GetTo() string {
	if er != nil && er.To != nil {
		return *er.To
	}
	return ""
}

func (er *ExchangeRate) SetTo(to *string) {
	er.To = to
}

func (er *ExchangeRate) GetRate() float64 {
	if er != nil && er.Rate != nil {
		return *er.Rate
	}
	return 0
}

func (er *ExchangeRate) SetRate(rate *float64) {
	er.Rate = rate
}

func (er *ExchangeRate) GetTimestamp() uint64 {
	if er != nil && er.Timestamp != nil {
		return *er.Timestamp
	}
	return 0
}

func (er *ExchangeRate) SetTimestamp(timestamp *uint64) {
	er.Timestamp = timestamp
}

func (er *ExchangeRate) GetCreateTime() uint64 {
	if er != nil && er.CreateTime != nil {
		return *er.CreateTime
	}
	return 0
}

func (er *ExchangeRate) SetCreateTime(createTime *uint64) {
	er.CreateTime = createTime
}
