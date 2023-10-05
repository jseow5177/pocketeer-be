package model

import (
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExchangeRate struct {
	ExchangeRateID primitive.ObjectID `bson:"_id,omitempty"`
	From           *string            `bson:"from,omitempty"`
	To             *string            `bson:"to,omitempty"`
	Rate           *float64           `bson:"rate,omitempty"`
	Timestamp      *uint64            `bson:"timestamp,omitempty"`
	CreateTime     *uint64            `bson:"create_time,omitempty"`
}

func ToExchangeRateModelFromEntity(er *entity.ExchangeRate) *ExchangeRate {
	if er == nil {
		return nil
	}

	objID := primitive.NilObjectID
	if primitive.IsValidObjectID(er.GetExchangeRateID()) {
		objID, _ = primitive.ObjectIDFromHex(er.GetExchangeRateID())
	}

	return &ExchangeRate{
		ExchangeRateID: objID,
		From:           er.From,
		To:             er.To,
		Rate:           er.Rate,
		Timestamp:      er.Timestamp,
		CreateTime:     er.CreateTime,
	}
}

func ToExchangeRateEntity(er *ExchangeRate) *entity.ExchangeRate {
	if er == nil {
		return nil
	}

	return entity.NewExchangeRate(
		er.GetFrom(),
		er.GetTo(),
		er.GetRate(),
		er.GetTimestamp(),
		entity.WithExchangeRateID(goutil.String(er.GetExchangeRateID())),
		entity.WithExchangeRateCreateTime(er.CreateTime),
	)
}

func (er *ExchangeRate) GetExchangeRateID() string {
	if er != nil {
		return er.ExchangeRateID.Hex()
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
