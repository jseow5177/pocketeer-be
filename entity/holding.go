package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

var (
	ErrSetCostAndValueForbidden = errors.New("set avg_cost or latest_value forbidden")
	ErrMustSetCostAndValue      = errors.New("avg_cost and latest_value must be set")
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

type HoldingUpdate struct {
	AvgCost     *float64
	LatestValue *float64
	UpdateTime  *uint64
}

func (hu *HoldingUpdate) GetAvgCost() float64 {
	if hu != nil && hu.AvgCost != nil {
		return *hu.AvgCost
	}
	return 0
}

func (hu *HoldingUpdate) GetLatestValue() float64 {
	if hu != nil && hu.LatestValue != nil {
		return *hu.LatestValue
	}
	return 0
}

func (hu *HoldingUpdate) GetUpdateTime() uint64 {
	if hu != nil && hu.UpdateTime != nil {
		return *hu.UpdateTime
	}
	return 0
}

func WithUpdateHoldingAvgCost(avgCost *float64) HoldingUpdateOption {
	return func(hu *HoldingUpdate) {
		hu.AvgCost = avgCost
	}
}

func WithUpdateHoldingLatestValue(latestValue *float64) HoldingUpdateOption {
	return func(hu *HoldingUpdate) {
		hu.LatestValue = latestValue
	}
}

type HoldingUpdateOption = func(hu *HoldingUpdate)

type Holding struct {
	UserID        *string
	HoldingID     *string
	AccountID     *string
	Symbol        *string
	HoldingStatus *uint32
	HoldingType   *uint32
	CreateTime    *uint64
	UpdateTime    *uint64
	TotalShares   *float64
	AvgCost       *float64
	LatestValue   *float64
	Quote         *Quote
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

func NewHolding(userID, accountID, symbol string, opts ...HoldingOption) (*Holding, error) {
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
	if err := h.checkOpts(); err != nil {
		return nil, err
	}
	return h, nil
}

func (h *Holding) checkOpts() error {
	if !h.IsCustom() {
		h.Symbol = goutil.String(strings.ToUpper(h.GetSymbol()))

		if h.AvgCost != nil || h.LatestValue != nil {
			return ErrSetCostAndValueForbidden
		}
	}

	if h.IsCustom() {
		if h.AvgCost == nil || h.LatestValue == nil {
			return ErrMustSetCostAndValue
		}
	}

	return nil
}

func (h *Holding) Update(hu *HoldingUpdate) (holdingUpdate *HoldingUpdate, hasUpdate bool, err error) {
	holdingUpdate = new(HoldingUpdate)

	if hu.AvgCost != nil && hu.GetAvgCost() != h.GetAvgCost() {
		hasUpdate = true
		h.AvgCost = hu.AvgCost
	}

	if hu.LatestValue != nil && hu.GetLatestValue() != h.GetLatestValue() {
		hasUpdate = true
		h.LatestValue = hu.LatestValue
	}

	if !hasUpdate {
		return
	}

	now := goutil.Uint64(uint64(time.Now().Unix()))
	h.UpdateTime = now

	if err = h.checkOpts(); err != nil {
		return nil, false, err
	}

	holdingUpdate.UpdateTime = now

	if hu.AvgCost != nil {
		holdingUpdate.AvgCost = h.AvgCost
	}

	if hu.LatestValue != nil {
		holdingUpdate.LatestValue = h.LatestValue
	}

	return
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
