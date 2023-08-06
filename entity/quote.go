package entity

import "github.com/jseow5177/pockteer-be/pkg/goutil"

type Currency string

const (
	CurrencyUSD Currency = "USD"
)

type Quote struct {
	LatestPrice   *float64
	Change        *float64 // LatestPrice - PreviousClose
	ChangePercent *float64
	PreviousClose *float64
	UpdateTime    *uint64
	Currency      *string
}

func (q *Quote) GetCurrency() string {
	if q != nil && q.Currency != nil {
		return *q.Currency
	}
	return ""
}

func (q *Quote) GetLatestPrice() float64 {
	if q != nil && q.LatestPrice != nil {
		return *q.LatestPrice
	}
	return 0
}

func (q *Quote) GetChange() float64 {
	if q != nil && q.Change != nil {
		return *q.Change
	}
	return 0
}

func (q *Quote) GetChangePercent() float64 {
	if q != nil && q.ChangePercent != nil {
		return *q.ChangePercent
	}
	return 0
}

func (q *Quote) GetPreviousClose() float64 {
	if q != nil && q.PreviousClose != nil {
		return *q.PreviousClose
	}
	return 0
}

func (q *Quote) GetUpdateTime() uint64 {
	if q != nil && q.UpdateTime != nil {
		return *q.UpdateTime
	}
	return 0
}

type QuoteOption = func(q *Quote)

func WithQuoteCurrency(currency *string) QuoteOption {
	return func(q *Quote) {
		q.Currency = currency
	}
}

func WithQuoteLatestPrice(latestPrice *float64) QuoteOption {
	return func(q *Quote) {
		q.LatestPrice = latestPrice
	}
}

func WithQuoteChange(change *float64) QuoteOption {
	return func(q *Quote) {
		q.Change = change
	}
}

func WithQuoteChangePercent(changePercent *float64) QuoteOption {
	return func(q *Quote) {
		q.ChangePercent = changePercent
	}
}

func WithQuotePreviousClose(previousClose *float64) QuoteOption {
	return func(q *Quote) {
		q.PreviousClose = previousClose
	}
}

func WithQuoteUpdateTime(updateTime *uint64) QuoteOption {
	return func(q *Quote) {
		q.UpdateTime = updateTime
	}
}

func NewQuote(opts ...QuoteOption) *Quote {
	q := &Quote{
		Change:        goutil.Float64(0),
		ChangePercent: goutil.Float64(0),
		LatestPrice:   goutil.Float64(0),
		PreviousClose: goutil.Float64(0),
		UpdateTime:    goutil.Uint64(0),
		Currency:      goutil.String(string(CurrencyUSD)),
	}
	for _, opt := range opts {
		opt(q)
	}
	return q
}
