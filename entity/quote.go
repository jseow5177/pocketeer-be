package entity

import "github.com/jseow5177/pockteer-be/pkg/goutil"

type Quote struct {
	QuoteID       *string
	Symbol        *string
	LatestPrice   *float64
	Change        *float64 // LatestPrice - PreviousClose
	ChangePercent *float64
	PreviousClose *float64
	UpdateTime    *uint64
	Currency      *string
}

func (q *Quote) GetQuoteID() string {
	if q != nil && q.QuoteID != nil {
		return *q.QuoteID
	}
	return ""
}

func (q *Quote) SetQuoteID(quoteID *string) {
	q.QuoteID = quoteID
}

func (q *Quote) GetSymbol() string {
	if q != nil && q.Symbol != nil {
		return *q.Symbol
	}
	return ""
}

func (q *Quote) SetSymbol(symbol *string) {
	q.Symbol = symbol
}

func (q *Quote) GetCurrency() string {
	if q != nil && q.Currency != nil {
		return *q.Currency
	}
	return ""
}

func (q *Quote) SetCurrency(currency *string) {
	q.Currency = currency
}

func (q *Quote) GetLatestPrice() float64 {
	if q != nil && q.LatestPrice != nil {
		return *q.LatestPrice
	}
	return 0
}

func (q *Quote) SetLatestPrice(latestPrice *float64) {
	q.LatestPrice = latestPrice
}

func (q *Quote) GetChange() float64 {
	if q != nil && q.Change != nil {
		return *q.Change
	}
	return 0
}

func (q *Quote) SetChange(change *float64) {
	q.Change = change
}

func (q *Quote) GetChangePercent() float64 {
	if q != nil && q.ChangePercent != nil {
		return *q.ChangePercent
	}
	return 0
}

func (q *Quote) SetChangePercent(changePercent *float64) {
	q.ChangePercent = changePercent
}

func (q *Quote) GetPreviousClose() float64 {
	if q != nil && q.PreviousClose != nil {
		return *q.PreviousClose
	}
	return 0
}

func (q *Quote) SetPreviousClose(previousClose *float64) {
	q.PreviousClose = previousClose
}

func (q *Quote) GetUpdateTime() uint64 {
	if q != nil && q.UpdateTime != nil {
		return *q.UpdateTime
	}
	return 0
}

func (q *Quote) SetUpdateTime(updateTime *uint64) {
	q.UpdateTime = updateTime
}

type QuoteOption = func(q *Quote)

func WithQuoteID(quoteID *string) QuoteOption {
	return func(q *Quote) {
		q.SetQuoteID(quoteID)
	}
}

func WithQuoteCurrency(currency *string) QuoteOption {
	return func(q *Quote) {
		q.SetCurrency(currency)
	}
}

func WithQuoteLatestPrice(latestPrice *float64) QuoteOption {
	return func(q *Quote) {
		q.SetLatestPrice(latestPrice)
	}
}

func WithQuoteChange(change *float64) QuoteOption {
	return func(q *Quote) {
		q.SetChange(change)
	}
}

func WithQuoteChangePercent(changePercent *float64) QuoteOption {
	return func(q *Quote) {
		q.SetChangePercent(changePercent)
	}
}

func WithQuotePreviousClose(previousClose *float64) QuoteOption {
	return func(q *Quote) {
		q.SetPreviousClose(previousClose)
	}
}

func WithQuoteUpdateTime(updateTime *uint64) QuoteOption {
	return func(q *Quote) {
		q.SetUpdateTime(updateTime)
	}
}

func NewQuote(symbol string, opts ...QuoteOption) *Quote {
	q := &Quote{
		QuoteID:       goutil.String(""),
		Symbol:        goutil.String(symbol),
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
