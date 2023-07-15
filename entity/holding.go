package entity

import (
	"strings"
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
	HoldingTypeInvalid HoldingType = iota
	HoldingTypeDefault
	HoldingTypeCustom
)

var HoldingTypes = map[uint32]string{
	uint32(HoldingTypeDefault): "default",
	uint32(HoldingTypeCustom):  "custom",
}

type Holding struct {
	UserID        *string
	HoldingID     *string
	AccountID     *string
	Symbol        *string
	HoldingStatus *uint32
	HoldingType   *uint32
	CreateTime    *uint64
	UpdateTime    *uint64

	// computed in real time
	TotalShares *float64
	AvgCost     *float64
	LatestValue *float64
	Quote       *Quote
}

type HoldingOption = func(h *Holding)

func WithHoldingID(holdingID *string) HoldingOption {
	return func(h *Holding) {
		h.HoldingID = holdingID
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

func WithHoldingTotalShares(totalShares *float64) HoldingOption {
	return func(h *Holding) {
		h.TotalShares = totalShares
	}
}

func WithHoldingAvgCost(wac *float64) HoldingOption {
	return func(h *Holding) {
		h.AvgCost = wac
	}
}

func WithHoldingLatestValue(latestValue *float64) HoldingOption {
	return func(h *Holding) {
		h.LatestValue = latestValue
	}
}

func WithQuote(quote *Quote) HoldingOption {
	return func(h *Holding) {
		h.Quote = quote
	}
}

func NewHolding(userID, accountID, symbol string, opts ...HoldingOption) *Holding {
	now := uint64(time.Now().Unix())
	h := &Holding{
		UserID:        goutil.String(userID),
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

func (h *Holding) checkOpts() {
	if !h.IsCustom() {
		h.Symbol = goutil.String(strings.ToUpper(h.GetSymbol()))
	}
}

func (h *Holding) GetHoldingID() string {
	if h != nil && h.HoldingID != nil {
		return *h.HoldingID
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

func (h *Holding) SetHoldingID(holdingID *string) {
	h.HoldingID = holdingID
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

func (h *Holding) GetTotalShares() float64 {
	if h != nil && h.TotalShares != nil {
		return *h.TotalShares
	}
	return 0
}

func (h *Holding) SetTotalShares(totalShares *float64) {
	h.TotalShares = totalShares
}

func (h *Holding) GetAvgCost() float64 {
	if h != nil && h.AvgCost != nil {
		return *h.AvgCost
	}
	return 0
}

func (h *Holding) SetAvgCost(avgCost *float64) {
	h.AvgCost = avgCost
}

func (h *Holding) GetLatestValue() float64 {
	if h != nil && h.LatestValue != nil {
		return *h.LatestValue
	}
	return 0
}

func (h *Holding) SetLatestValue(latestValue *float64) {
	h.LatestValue = latestValue
}

func (h *Holding) GetQuote() *Quote {
	if h != nil && h.Quote != nil {
		return h.Quote
	}
	return nil
}

func (h *Holding) SetQuote(quote *Quote) {
	h.Quote = quote
}

func (h *Holding) IsCustom() bool {
	return h.GetHoldingType() == uint32(HoldingTypeCustom)
}
