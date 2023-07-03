package entity

import (
	"time"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

type HoldingStatus uint32

const (
	HoldingStatusInvalid AccountStatus = iota
	HoldingStatusNormal
	HoldingStatusDeleted
)

type HoldingType uint32

const (
	HoldingTypeDefault HoldingType = iota
	HoldingTypeCustom
)

var HoldingTypes = map[uint32]string{
	uint32(HoldingTypeDefault): "default",
	uint32(HoldingTypeCustom):  "custom",
}

type Holding struct {
	HoldingID     *string
	AccountID     *string
	Symbol        *string
	HoldingStatus *uint32
	HoldingType   *uint32
	CreateTime    *uint64
	UpdateTime    *uint64
}

type HoldingOption = func(h *Holding)

func WithHoldingID(holdingID *string) HoldingOption {
	return func(h *Holding) {
		h.HoldingID = holdingID
	}
}

func WithHoldingAccountID(holdingID *string) HoldingOption {
	return func(h *Holding) {
		h.AccountID = holdingID
	}
}

func WithSymbol(symbol *string) HoldingOption {
	return func(h *Holding) {
		h.Symbol = symbol
	}
}

func WithHoldingStatus(holdingStatus *uint32) HoldingOption {
	return func(h *Holding) {
		h.HoldingStatus = holdingStatus
	}
}

func WithHoldingType(holdingType *uint32) HoldingOption {
	return func(h *Holding) {
		h.HoldingType = holdingType
	}
}

func WithHoldingCreateTime(createTime *uint64) HoldingOption {
	return func(h *Holding) {
		h.CreateTime = createTime
	}
}

func WithHoldingUpdateTime(updateTime *uint64) HoldingOption {
	return func(h *Holding) {
		h.UpdateTime = updateTime
	}
}

func NewHolding(accountID, symbol string, opts ...HoldingOption) *Holding {
	now := uint64(time.Now().Unix())
	h := &Holding{
		AccountID:     goutil.String(accountID),
		Symbol:        goutil.String(symbol),
		HoldingStatus: goutil.Uint32(uint32(HoldingStatusNormal)),
		HoldingType:   goutil.Uint32(uint32(HoldingTypeCustom)),
		CreateTime:    goutil.Uint64(now),
		UpdateTime:    goutil.Uint64(now),
	}
	for _, opt := range opts {
		opt(h)
	}
	h.checkOpts()
	return h
}

func setHolding(h *Holding, opts ...HoldingOption) {
	if h == nil {
		return
	}

	for _, opt := range opts {
		opt(h)
	}
}

func (h *Holding) checkOpts() {}

func (h *Holding) GetHoldingID() string {
	if h != nil && h.HoldingID != nil {
		return *h.HoldingID
	}
	return ""
}

func (h *Holding) GetAccountID() string {
	if h != nil && h.AccountID != nil {
		return *h.AccountID
	}
	return ""
}

func (h *Holding) SetHoldingID(holdingID string) {
	setHolding(h, WithHoldingID(goutil.String(holdingID)))
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
