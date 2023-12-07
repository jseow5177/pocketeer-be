package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/util"
)

var (
	ErrSetCostValueForbidden = errors.New("set total_cost or latest_value forbidden")
	ErrMustSetCostAndValue   = errors.New("total_cost and latest_value must be set")
	ErrHoldingCannotHaveLots = errors.New("holding cannot have lots")
	ErrMismatchCurrency      = errors.New("mismatch currency")
	ErrCannotChangeSymbol    = errors.New("cannot change symbol")
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

type HoldingUpdateOption func(h *Holding)

func WithUpdateHoldingTotalCost(totalCost *float64) HoldingUpdateOption {
	return func(h *Holding) {
		if totalCost != nil {
			h.SetTotalCost(totalCost)
		}
	}
}

func WithUpdateHoldingLatestValue(latestValue *float64) HoldingUpdateOption {
	return func(h *Holding) {
		if latestValue != nil {
			h.SetLatestValue(latestValue)
		}
	}
}

func WithUpdateHoldingSymbol(symbol *string) HoldingUpdateOption {
	return func(h *Holding) {
		if symbol != nil {
			h.SetSymbol(symbol)
		}
	}
}

func WithUpdateHoldingStatus(holdingStatus *uint32) HoldingUpdateOption {
	return func(h *Holding) {
		if holdingStatus != nil {
			h.SetHoldingStatus(holdingStatus)
		}
	}
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
	Currency      *string
	TotalCost     *float64 // stored for custom, computed for default
	LatestValue   *float64 // stored for custom, computed for default

	TotalShares     *float64 // no-op for custom, computed for default from lots
	AvgCostPerShare *float64 // no-op for custom, computed for default from lots
	Quote           *Quote   // no-op for custom, computed for default from external API
	Gain            *float64 // no-op for customm, computed for default from lots
	PercentGain     *float64 // no-op for customm, computed for default from lots

	Lots []*Lot
}

type HoldingOption func(h *Holding)

func WithHoldingID(holdingID *string) HoldingOption {
	return func(h *Holding) {
		if holdingID != nil {
			h.SetHoldingID(holdingID)
		}
	}
}

func WithHoldingStatus(holdingStatus *uint32) HoldingOption {
	return func(h *Holding) {
		if holdingStatus != nil {
			h.SetHoldingStatus(holdingStatus)
		}
	}
}

func WithHoldingType(holdingType *uint32) HoldingOption {
	return func(h *Holding) {
		if holdingType != nil {
			h.SetHoldingType(holdingType)
		}
	}
}

func WithHoldingCreateTime(createTime *uint64) HoldingOption {
	return func(h *Holding) {
		if createTime != nil {
			h.SetCreateTime(createTime)
		}
	}
}

func WithHoldingUpdateTime(updateTime *uint64) HoldingOption {
	return func(h *Holding) {
		if updateTime != nil {
			h.SetUpdateTime(updateTime)
		}
	}
}

func WithHoldingCurrency(currency *string) HoldingOption {
	return func(h *Holding) {
		if currency != nil {
			h.SetCurrency(currency)
		}
	}
}

func WithHoldingTotalCost(totalCost *float64) HoldingOption {
	return func(h *Holding) {
		if totalCost != nil {
			h.SetTotalCost(totalCost)
		}
	}
}

func WithHoldingLatestValue(latestValue *float64) HoldingOption {
	return func(h *Holding) {
		if latestValue != nil {
			h.SetLatestValue(latestValue)
		}
	}
}

func WithHoldingTotalShares(totalShares *float64) HoldingOption {
	return func(h *Holding) {
		if totalShares != nil {
			h.SetTotalShares(totalShares)
		}
	}
}

func WithHoldingLots(lots []*Lot) HoldingOption {
	return func(h *Holding) {
		if lots != nil {
			h.SetLots(lots)
		}
	}
}

func (h *Holding) Clone() (*Holding, error) {
	return NewHolding(
		h.GetUserID(),
		h.GetAccountID(),
		h.GetSymbol(),
		WithHoldingID(goutil.String(h.GetHoldingID())),
		WithHoldingType(h.HoldingType),
		WithHoldingStatus(h.HoldingStatus),
		WithHoldingCreateTime(h.CreateTime),
		WithHoldingUpdateTime(h.UpdateTime),
		WithHoldingTotalCost(h.TotalCost),
		WithHoldingLatestValue(h.LatestValue),
		WithHoldingCurrency(h.Currency),
	)
}

func NewHolding(userID, accountID, symbol string, opts ...HoldingOption) (*Holding, error) {
	now := uint64(time.Now().UnixMilli())
	h := &Holding{
		UserID:        goutil.String(userID),
		AccountID:     goutil.String(accountID),
		Symbol:        goutil.String(symbol),
		HoldingStatus: goutil.Uint32(uint32(HoldingStatusNormal)),
		HoldingType:   goutil.Uint32(uint32(HoldingTypeCustom)),
		CreateTime:    goutil.Uint64(now),
		UpdateTime:    goutil.Uint64(now),
		Currency:      goutil.String(string(CurrencySGD)),
	}

	for _, opt := range opts {
		opt(h)
	}

	if err := h.validate(); err != nil {
		return nil, err
	}

	return h, nil
}

func (h *Holding) validate() error {
	if h.IsDefault() {
		h.Symbol = goutil.String(strings.ToUpper(h.GetSymbol()))

		// for non-custom type, set cannot cost, value, and shares
		if h.TotalCost != nil || h.LatestValue != nil {
			return ErrSetCostValueForbidden
		}
	}

	if h.IsCustom() {
		if h.TotalCost == nil || h.LatestValue == nil {
			return ErrMustSetCostAndValue
		}

		if len(h.Lots) > 0 {
			return ErrHoldingCannotHaveLots
		}
	}

	return nil
}

type HoldingUpdate struct {
	TotalCost     *float64
	LatestValue   *float64
	Symbol        *string
	UpdateTime    *uint64
	HoldingStatus *uint32
}

func (h *Holding) ToHoldingUpdate(old *Holding) *HoldingUpdate {
	var (
		hasUpdate bool

		hu = &HoldingUpdate{
			UpdateTime: h.UpdateTime,
		}
	)

	if old.GetTotalCost() != h.GetTotalCost() {
		hasUpdate = true
		hu.TotalCost = h.TotalCost
	}

	if old.GetLatestValue() != h.GetLatestValue() {
		hasUpdate = true
		hu.LatestValue = h.LatestValue
	}

	if old.GetSymbol() != h.GetSymbol() {
		hasUpdate = true
		hu.Symbol = h.Symbol
	}

	if old.GetHoldingStatus() != h.GetHoldingStatus() {
		hasUpdate = true
		hu.HoldingStatus = h.HoldingStatus
	}

	if hasUpdate {
		return hu
	}

	return nil
}

func (h *Holding) Update(hus ...HoldingUpdateOption) (*HoldingUpdate, error) {
	if len(hus) == 0 {
		return nil, nil
	}

	old, err := h.Clone()
	if err != nil {
		return nil, err
	}

	for _, hu := range hus {
		hu(h)
	}

	// check
	if err := h.validate(); err != nil {
		return nil, err
	}

	now := goutil.Uint64(uint64(time.Now().UnixMilli()))
	h.SetUpdateTime(now)

	return h.ToHoldingUpdate(old), nil
}

func (h *Holding) GetHoldingID() string {
	if h != nil && h.HoldingID != nil {
		return *h.HoldingID
	}
	return ""
}

func (h *Holding) SetHoldingID(holdingID *string) {
	h.HoldingID = holdingID
}

func (h *Holding) GetUserID() string {
	if h != nil && h.UserID != nil {
		return *h.UserID
	}
	return ""
}

func (h *Holding) SetUserID(userID *string) {
	h.UserID = userID
}

func (h *Holding) GetAccountID() string {
	if h != nil && h.AccountID != nil {
		return *h.AccountID
	}
	return ""
}

func (h *Holding) SetAccountID(accountID *string) {
	h.AccountID = accountID
}

func (h *Holding) GetSymbol() string {
	if h != nil && h.Symbol != nil {
		return *h.Symbol
	}
	return ""
}

func (h *Holding) SetSymbol(symbol *string) {
	h.Symbol = symbol
}

func (h *Holding) GetHoldingStatus() uint32 {
	if h != nil && h.HoldingStatus != nil {
		return *h.HoldingStatus
	}
	return 0
}

func (h *Holding) SetHoldingStatus(holdingStatus *uint32) {
	h.HoldingStatus = holdingStatus
}

func (h *Holding) GetHoldingType() uint32 {
	if h != nil && h.HoldingType != nil {
		return *h.HoldingType
	}
	return 0
}

func (h *Holding) SetHoldingType(holdingType *uint32) {
	h.HoldingType = holdingType
}

func (h *Holding) GetCreateTime() uint64 {
	if h != nil && h.CreateTime != nil {
		return *h.CreateTime
	}
	return 0
}

func (h *Holding) SetCreateTime(createTime *uint64) {
	h.CreateTime = createTime
}

func (h *Holding) GetUpdateTime() uint64 {
	if h != nil && h.UpdateTime != nil {
		return *h.UpdateTime
	}
	return 0
}

func (h *Holding) SetUpdateTime(updateTime *uint64) {
	h.UpdateTime = updateTime
}

func (h *Holding) GetTotalShares() float64 {
	if h != nil && h.TotalShares != nil {
		return *h.TotalShares
	}
	return 0
}

func (h *Holding) SetTotalShares(totalShares *float64) {
	h.TotalShares = totalShares

	if totalShares != nil {
		ts := util.RoundFloatToStandardDP(*totalShares)
		h.TotalShares = goutil.Float64(ts)
	}
}

func (h *Holding) GetAvgCostPerShare() float64 {
	if h != nil && h.AvgCostPerShare != nil {
		return *h.AvgCostPerShare
	}
	return 0
}

func (h *Holding) SetAvgCostPerShare(avgCostPerShare *float64) {
	h.AvgCostPerShare = avgCostPerShare

	if avgCostPerShare != nil {
		acps := util.RoundFloatToStandardDP(*avgCostPerShare)
		h.AvgCostPerShare = goutil.Float64(acps)
	}
}

func (h *Holding) GetTotalCost() float64 {
	if h != nil && h.TotalCost != nil {
		return *h.TotalCost
	}
	return 0
}

func (h *Holding) SetTotalCost(totalCost *float64) {
	h.TotalCost = totalCost

	if totalCost != nil {
		tc := util.RoundFloatToStandardDP(*totalCost)
		h.TotalCost = goutil.Float64(tc)
	}
}

func (h *Holding) GetLatestValue() float64 {
	if h != nil && h.LatestValue != nil {
		return *h.LatestValue
	}
	return 0
}

func (h *Holding) SetLatestValue(latestValue *float64) {
	h.LatestValue = latestValue

	if latestValue != nil {
		lv := util.RoundFloatToStandardDP(*latestValue)
		h.LatestValue = goutil.Float64(lv)
	}
}

func (h *Holding) GetCurrency() string {
	if h != nil && h.Currency != nil {
		return *h.Currency
	}
	return ""
}

func (h *Holding) SetCurrency(currency *string) {
	h.Currency = currency
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

func (h *Holding) GetGain() float64 {
	if h != nil && h.Gain != nil {
		return *h.Gain
	}
	return 0
}

func (h *Holding) SetGain(gain *float64) {
	h.Gain = gain

	if gain != nil {
		g := util.RoundFloatToStandardDP(*gain)
		h.Gain = goutil.Float64(g)
	}
}

func (h *Holding) GetPercentGain() float64 {
	if h != nil && h.PercentGain != nil {
		return *h.PercentGain
	}
	return 0
}

func (h *Holding) SetPercentGain(percentGain *float64) {
	h.PercentGain = percentGain

	if percentGain != nil {
		pg := util.RoundFloatToStandardDP(*percentGain)
		h.PercentGain = goutil.Float64(pg)
	}
}

func (h *Holding) GetLots() []*Lot {
	if h != nil && h.Lots != nil {
		return h.Lots
	}
	return nil
}

func (h *Holding) SetLots(ls []*Lot) {
	h.Lots = ls
}

func (h *Holding) IsCustom() bool {
	return h.GetHoldingType() == uint32(HoldingTypeCustom)
}

func (h *Holding) IsDefault() bool {
	return h.GetHoldingType() == uint32(HoldingTypeDefault)
}

func (h *Holding) CanHaveLots() bool {
	return h.IsDefault()
}

// Compute the latest value, total cost, avg cost, gain, and percent gain of a holding.
//
// No currency conversion is needed as holding, lots, and security currency should be same.
func (h *Holding) ComputeCostGainAndValue() {
	if !h.IsDefault() {
		gain := h.GetLatestValue() - h.GetTotalCost()
		h.SetGain(goutil.Float64(gain))

		var percentGain *float64
		if h.GetTotalCost() > 0 {
			percentGain = goutil.Float64(gain * 100 / h.GetTotalCost())
		}
		h.SetPercentGain(percentGain)
		return
	}

	var (
		totalCost   float64
		totalShares float64
	)
	for _, l := range h.Lots {
		totalCost += l.GetCostPerShare() * l.GetShares()
		totalShares += l.GetShares()
	}

	var avgCostPerShare float64
	if totalShares > 0 {
		avgCostPerShare = totalCost / totalShares
	}

	latestValue := totalShares * h.Quote.GetLatestPrice()

	h.SetTotalShares(goutil.Float64(totalShares))
	h.SetTotalCost(goutil.Float64(totalCost))
	h.SetAvgCostPerShare(goutil.Float64(avgCostPerShare))
	h.SetLatestValue(goutil.Float64(latestValue))

	gain := latestValue - totalCost
	h.SetGain(goutil.Float64(gain))

	var percentGain *float64
	if totalCost > 0 {
		percentGain = goutil.Float64(gain * 100 / totalCost)
	}
	h.SetPercentGain(percentGain)
}
