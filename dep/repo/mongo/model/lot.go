package model

import (
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Lot struct {
	LotID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID       *string            `bson:"user_id,omitempty"`
	HoldingID    *string            `bson:"holding_id,omitempty"`
	Shares       *float64           `bson:"shares,omitempty"`
	CostPerShare *float64           `bson:"cost_per_share,omitempty"`
	LotStatus    *uint32            `bson:"lot_status,omitempty"`
	TradeDate    *uint64            `bson:"trade_date,omitempty"`
	CreateTime   *uint64            `bson:"create_time,omitempty"`
	UpdateTime   *uint64            `bson:"update_time,omitempty"`
	Currency     *string            `bson:"currency,omitempty"`
}

func ToLotModelFromEntity(l *entity.Lot) *Lot {
	if l == nil {
		return nil
	}

	objID := primitive.NilObjectID
	if primitive.IsValidObjectID(l.GetLotID()) {
		objID, _ = primitive.ObjectIDFromHex(l.GetLotID())
	}

	return &Lot{
		LotID:        objID,
		UserID:       l.UserID,
		HoldingID:    l.HoldingID,
		Shares:       l.Shares,
		CostPerShare: l.CostPerShare,
		LotStatus:    l.LotStatus,
		TradeDate:    l.TradeDate,
		CreateTime:   l.CreateTime,
		UpdateTime:   l.UpdateTime,
		Currency:     l.Currency,
	}
}

func ToLotModelFromUpdate(lu *entity.LotUpdate) *Lot {
	if lu == nil {
		return nil
	}

	return &Lot{
		Shares:       lu.Shares,
		CostPerShare: lu.CostPerShare,
		TradeDate:    lu.TradeDate,
		UpdateTime:   lu.UpdateTime,
		LotStatus:    lu.LotStatus,
	}
}

func ToLotEntity(l *Lot) *entity.Lot {
	if l == nil {
		return nil
	}

	return entity.NewLot(
		l.GetUserID(),
		l.GetHoldingID(),
		entity.WithLotID(goutil.String(l.GetLotID())),
		entity.WithLotShares(l.Shares),
		entity.WithLotStatus(l.LotStatus),
		entity.WithLotCostPerShare(l.CostPerShare),
		entity.WithLotTradeDate(l.TradeDate),
		entity.WithLotCreateTime(l.CreateTime),
		entity.WithLotUpdateTime(l.UpdateTime),
		entity.WithLotCurrency(l.Currency),
	)
}

func (l *Lot) GetLotID() string {
	if l != nil {
		return l.LotID.Hex()
	}
	return ""
}

func (l *Lot) GetUserID() string {
	if l != nil && l.UserID != nil {
		return *l.UserID
	}
	return ""
}

func (l *Lot) GetHoldingID() string {
	if l != nil && l.HoldingID != nil {
		return *l.HoldingID
	}
	return ""
}

func (l *Lot) GetShares() float64 {
	if l != nil && l.Shares != nil {
		return *l.Shares
	}
	return 0
}

func (l *Lot) GetCostPerShare() float64 {
	if l != nil && l.CostPerShare != nil {
		return *l.CostPerShare
	}
	return 0
}

func (l *Lot) GetTradeDate() uint64 {
	if l != nil && l.TradeDate != nil {
		return *l.TradeDate
	}
	return 0
}

func (l *Lot) GetLotStatus() uint32 {
	if l != nil && l.LotStatus != nil {
		return *l.LotStatus
	}
	return 0
}

func (l *Lot) GetCreateTime() uint64 {
	if l != nil && l.CreateTime != nil {
		return *l.CreateTime
	}
	return 0
}

func (l *Lot) GetUpdateTime() uint64 {
	if l != nil && l.UpdateTime != nil {
		return *l.UpdateTime
	}
	return 0
}

func (l *Lot) GetCurrency() string {
	if l != nil && l.Currency != nil {
		return *l.Currency
	}
	return ""
}
