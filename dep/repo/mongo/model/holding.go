package model

import (
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Holding struct {
	HoldingID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID        *string            `bson:"user_id,omitempty"`
	AccountID     *string            `bson:"account_id,omitempty"`
	Symbol        *string            `bson:"symbol,omitempty"`
	HoldingStatus *uint32            `bson:"holding_status,omitempty"`
	HoldingType   *uint32            `bson:"holding_type,omitempty"`
	CreateTime    *uint64            `bson:"create_time,omitempty"`
	UpdateTime    *uint64            `bson:"update_time,omitempty"`

	// for custom holding
	TotalCost   *float64 `bson:"total_cost,omitempty"`
	LatestValue *float64 `bson:"latest_value,omitempty"`
}

func ToHoldingModelFromEntity(h *entity.Holding) *Holding {
	if h == nil {
		return nil
	}

	objID := primitive.NilObjectID
	if primitive.IsValidObjectID(h.GetHoldingID()) {
		objID, _ = primitive.ObjectIDFromHex(h.GetHoldingID())
	}

	return &Holding{
		HoldingID:     objID,
		UserID:        h.UserID,
		AccountID:     h.AccountID,
		Symbol:        h.Symbol,
		HoldingType:   h.HoldingType,
		HoldingStatus: h.HoldingStatus,
		CreateTime:    h.CreateTime,
		UpdateTime:    h.UpdateTime,
		TotalCost:     h.TotalCost,
		LatestValue:   h.LatestValue,
	}
}

func ToHoldingModelFromUpdate(hu *entity.HoldingUpdate) *Holding {
	if hu == nil {
		return nil
	}

	return &Holding{
		TotalCost:   hu.TotalCost,
		LatestValue: hu.LatestValue,
	}
}

func ToHoldingEntity(h *Holding) (*entity.Holding, error) {
	if h == nil {
		return nil, nil
	}

	return entity.NewHolding(
		h.GetUserID(),
		h.GetAccountID(),
		h.GetSymbol(),
		entity.WithHoldingID(goutil.String(h.GetHoldingID())),
		entity.WithHoldingType(h.HoldingType),
		entity.WithHoldingStatus(h.HoldingStatus),
		entity.WithHoldingCreateTime(h.CreateTime),
		entity.WithHoldingUpdateTime(h.UpdateTime),
		entity.WithHoldingTotalCost(h.TotalCost),
		entity.WithHoldingLatestValue(h.LatestValue),
	)
}

func (h *Holding) GetHoldingID() string {
	if h != nil {
		return h.HoldingID.Hex()
	}
	return ""
}

func (h *Holding) GetUserID() string {
	if h != nil && h.UserID != nil {
		return *h.UserID
	}
	return ""
}

func (h *Holding) GetAccountID() string {
	if h != nil && h.AccountID != nil {
		return *h.AccountID
	}
	return ""
}

func (h *Holding) GetSymbol() string {
	if h != nil && h.Symbol != nil {
		return *h.Symbol
	}
	return ""
}

func (h *Holding) GetHoldingStatus() uint32 {
	if h != nil && h.HoldingStatus != nil {
		return *h.HoldingStatus
	}
	return 0
}

func (h *Holding) GetHoldingType() uint32 {
	if h != nil && h.HoldingType != nil {
		return *h.HoldingType
	}
	return 0
}

func (h *Holding) GetCreateTime() uint64 {
	if h != nil && h.CreateTime != nil {
		return *h.CreateTime
	}
	return 0
}

func (h *Holding) GetUpdateTime() uint64 {
	if h != nil && h.UpdateTime != nil {
		return *h.UpdateTime
	}
	return 0
}

func (h *Holding) GetTotalCost() float64 {
	if h != nil && h.TotalCost != nil {
		return *h.TotalCost
	}
	return 0
}

func (h *Holding) GetLatestValue() float64 {
	if h != nil && h.LatestValue != nil {
		return *h.LatestValue
	}
	return 0
}
