package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/util"
)

var (
	ErrSetCostValueSharesForbidden = errors.New("set total_cost or latest_value or total_shares forbidden")
	ErrMustSetCostAndValue         = errors.New("total_cost and latest_value must be set")
	ErrHoldingCannotHaveLots       = errors.New("holding cannot have lots")
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
	TotalCost     *float64
	LatestValue   *float64
	Symbol        *string
	UpdateTime    *uint64
	HoldingStatus *uint32
}

func (hu *HoldingUpdate) GetTotalCost() float64 {
	if hu != nil && hu.TotalCost != nil {
		return *hu.TotalCost
	}
	return 0
}

func (hu *HoldingUpdate) SetTotalCost(totalCost *float64) {
	hu.TotalCost = totalCost

	if totalCost != nil {
		tc := util.RoundFloatToStandardDP(*totalCost)
		hu.TotalCost = goutil.Float64(tc)
	}
}

func (hu *HoldingUpdate) GetLatestValue() float64 {
	if hu != nil && hu.LatestValue != nil {
		return *hu.LatestValue
	}
	return 0
}

func (hu *HoldingUpdate) SetLatestValue(latestValue *float64) {
	hu.LatestValue = latestValue

	if latestValue != nil {
		lv := util.RoundFloatToStandardDP(*latestValue)
		hu.LatestValue = goutil.Float64(lv)
	}
}

func (hu *HoldingUpdate) GetSymbol() string {
	if hu != nil && hu.Symbol != nil {
		return *hu.Symbol
	}
	return ""
}

func (hu *HoldingUpdate) SetSymbol(symbol *string) {
	hu.Symbol = symbol
}

func (hu *HoldingUpdate) GetUpdateTime() uint64 {
	if hu != nil && hu.UpdateTime != nil {
		return *hu.UpdateTime
	}
	return 0
}

func (hu *HoldingUpdate) SetUpdateTime(updateTime *uint64) {
	hu.UpdateTime = updateTime
}

func (hu *HoldingUpdate) GetHoldingStatus() uint32 {
	if hu != nil && hu.HoldingStatus != nil {
		return *hu.HoldingStatus
	}
	return 0
}

func (hu *HoldingUpdate) SetHoldingStatus(holdingStatus *uint32) {
	hu.HoldingStatus = holdingStatus
}

func WithUpdateHoldingTotalCost(totalCost *float64) HoldingUpdateOption {
	return func(hu *HoldingUpdate) {
		hu.SetTotalCost(totalCost)
	}
}

func WithUpdateHoldingLatestValue(latestValue *float64) HoldingUpdateOption {
	return func(hu *HoldingUpdate) {
		hu.SetLatestValue(latestValue)
	}
}

func WithUpdateHoldingSymbol(symbol *string) HoldingUpdateOption {
	return func(hu *HoldingUpdate) {
		hu.SetSymbol(symbol)
	}
}

func WithUpdateHoldingStatus(holdingStatus *uint32) HoldingUpdateOption {
	return func(hu *HoldingUpdate) {
		hu.SetHoldingStatus(holdingStatus)
	}
}

type HoldingUpdateOption = func(hu *HoldingUpdate)

func NewHoldingUpdate(opts ...HoldingUpdateOption) *HoldingUpdate {
	hu := new(HoldingUpdate)
	for _, opt := range opts {
		opt(hu)
	}
	return hu
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

	TotalShares     *float64 // no-op for customm, computed for default from lots
	AvgCostPerShare *float64 // no-op for custom, computed for default from lots
	Quote           *Quote   // no-op for custom, computed for default from external API
	TotalCost       *float64 // stored for custom, computed for default
	LatestValue     *float64 // stored for custom, computed for default

	Lots []*Lot
}

type HoldingOption = func(h *Holding)

func WithHoldingID(holdingID *string) HoldingOption {
	return func(h *Holding) {
		h.SetHoldingID(holdingID)
	}
}

func WithHoldingStatus(holdingStatus *uint32) HoldingOption {
	return func(h *Holding) {
		h.SetHoldingStatus(holdingStatus)
	}
}

func WithHoldingType(holdingType *uint32) HoldingOption {
	return func(h *Holding) {
		h.SetHoldingType(holdingType)
	}
}

func WithHoldingCreateTime(createTime *uint64) HoldingOption {
	return func(h *Holding) {
		h.SetCreateTime(createTime)
	}
}

func WithHoldingUpdateTime(updateTime *uint64) HoldingOption {
	return func(h *Holding) {
		h.SetUpdateTime(updateTime)
	}
}

func WithHoldingTotalCost(totalCost *float64) HoldingOption {
	return func(h *Holding) {
		h.SetTotalCost(totalCost)
	}
}

func WithHoldingLatestValue(latestValue *float64) HoldingOption {
	return func(h *Holding) {
		h.SetLatestValue(latestValue)
	}
}

func WithHoldingTotalShares(totalShares *float64) HoldingOption {
	return func(h *Holding) {
		h.SetTotalShares(totalShares)
	}
}

func WithHoldingLots(lots []*Lot) HoldingOption {
	return func(h *Holding) {
		h.SetLots(lots)
	}
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
	if h.IsDefault() {
		h.Symbol = goutil.String(strings.ToUpper(h.GetSymbol()))

		// for non-custom type, set cannot cost, value, and shares
		if h.TotalCost != nil || h.LatestValue != nil || h.TotalShares != nil {
			return ErrSetCostValueSharesForbidden
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

func (h *Holding) Update(hu *HoldingUpdate) (*HoldingUpdate, error) {
	var (
		hasUpdate     bool
		holdingUpdate = new(HoldingUpdate)
	)

	if hu.TotalCost != nil && hu.GetTotalCost() != h.GetTotalCost() {
		hasUpdate = true
		h.SetTotalCost(hu.TotalCost)

		defer func() {
			holdingUpdate.SetTotalCost(h.TotalCost)
		}()
	}

	if hu.LatestValue != nil && hu.GetLatestValue() != h.GetLatestValue() {
		hasUpdate = true
		h.SetLatestValue(hu.LatestValue)

		defer func() {
			holdingUpdate.SetLatestValue(h.LatestValue)
		}()
	}

	if hu.Symbol != nil && hu.GetSymbol() != h.GetSymbol() {
		hasUpdate = true
		h.SetSymbol(hu.Symbol)

		defer func() {
			holdingUpdate.SetSymbol(h.Symbol)
		}()
	}

	if !hasUpdate {
		return nil, nil
	}

	now := goutil.Uint64(uint64(time.Now().UnixMilli()))
	h.SetUpdateTime(now)

	if err := h.checkOpts(); err != nil {
		return nil, err
	}

	holdingUpdate.SetUpdateTime(now)

	return holdingUpdate, nil
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

func (h *Holding) GetQuote() *Quote {
	if h != nil && h.Quote != nil {
		return h.Quote
	}
	return nil
}

func (h *Holding) SetQuote(quote *Quote) {
	h.Quote = quote
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

func (h *Holding) ComputeSharesCostAndValue() {
	if !h.IsDefault() {
		return
	}

	// TODO
	conv := config.USDToSGD

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
	h.SetTotalCost(goutil.Float64(totalCost * conv))
	h.SetAvgCostPerShare(goutil.Float64(avgCostPerShare * conv))
	h.SetLatestValue(goutil.Float64(latestValue * conv))
}
