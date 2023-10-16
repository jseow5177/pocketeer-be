package model

import (
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Quote struct {
	QuoteID       primitive.ObjectID `bson:"_id,omitempty"`
	Symbol        *string            `bson:"symbol,omitempty"`
	LatestPrice   *float64           `bson:"latest_price,omitempty"`
	Change        *float64           `bson:"change,omitempty"`
	ChangePercent *float64           `bson:"change_percent,omitempty"`
	PreviousClose *float64           `bson:"previous_close,omitempty"`
	UpdateTime    *uint64            `bson:"update_time,omitempty"`
}

func (q *Quote) GetQuoteID() string {
	if q != nil {
		return q.QuoteID.Hex()
	}
	return ""
}

func (q *Quote) GetSymbol() string {
	if q != nil && q.Symbol != nil {
		return *q.Symbol
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

func ToQuoteModelFromEntity(q *entity.Quote) *Quote {
	if q == nil {
		return nil
	}

	objID := primitive.NilObjectID
	if primitive.IsValidObjectID(q.GetQuoteID()) {
		objID, _ = primitive.ObjectIDFromHex(q.GetQuoteID())
	}

	return &Quote{
		QuoteID:       objID,
		Symbol:        q.Symbol,
		LatestPrice:   q.LatestPrice,
		Change:        q.Change,
		ChangePercent: q.ChangePercent,
		PreviousClose: q.PreviousClose,
		UpdateTime:    q.UpdateTime,
	}
}

func ToQuoteEntity(q *Quote) *entity.Quote {
	if q == nil {
		return nil
	}

	return entity.NewQuote(
		q.GetSymbol(),
		entity.WithQuoteID(goutil.String(q.GetQuoteID())),
		entity.WithQuoteLatestPrice(q.LatestPrice),
		entity.WithQuoteChange(q.Change),
		entity.WithQuoteChangePercent(q.ChangePercent),
		entity.WithQuotePreviousClose(q.PreviousClose),
		entity.WithQuoteUpdateTime(q.UpdateTime),
	)
}
