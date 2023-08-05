package model

import "github.com/jseow5177/pockteer-be/entity"

type Quote struct {
	LatestPrice   *float64
	Change        *float64
	ChangePercent *float64
	PreviousClose *float64
	UpdateTime    *uint64
	Currency      *string
}

func ToQuoteModelFromEntity(q *entity.Quote) *Quote {
	if q == nil {
		return nil
	}

	return &Quote{
		LatestPrice:   q.LatestPrice,
		Change:        q.Change,
		ChangePercent: q.ChangePercent,
		PreviousClose: q.PreviousClose,
		UpdateTime:    q.UpdateTime,
		Currency:      q.Currency,
	}
}

func ToQuoteEntity(q *Quote) *entity.Quote {
	if q == nil {
		return nil
	}

	return entity.NewQuote(
		entity.WithQuoteLatestPrice(q.LatestPrice),
		entity.WithQuoteChange(q.Change),
		entity.WithQuoteChangePercent(q.ChangePercent),
		entity.WithQuotePreviousClose(q.PreviousClose),
		entity.WithQuoteCurrency(q.Currency),
		entity.WithQuoteUpdateTime(q.UpdateTime),
	)
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
