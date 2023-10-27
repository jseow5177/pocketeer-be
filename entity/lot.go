package entity

import (
	"time"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/util"
)

type LotStatus uint32

const (
	LotStatusInvalid AccountStatus = iota
	LotStatusNormal
	LotStatusDeleted
)

type LotUpdateOption func(l *Lot)

func WithUpdateLotShares(shares *float64) LotUpdateOption {
	return func(l *Lot) {
		if shares != nil {
			l.SetShares(shares)
		}
	}
}

func WithUpdateLotCostPerShare(costPerShare *float64) LotUpdateOption {
	return func(l *Lot) {
		if costPerShare != nil {
			l.SetCostPerShare(costPerShare)
		}
	}
}

func WithUpdateLotTradeDate(tradeDate *uint64) LotUpdateOption {
	return func(l *Lot) {
		if tradeDate != nil {
			l.SetTradeDate(tradeDate)
		}
	}
}

func WithUpdateLotStatus(lotStatus *uint32) LotUpdateOption {
	return func(l *Lot) {
		if lotStatus != nil {
			l.SetLotStatus(lotStatus)
		}
	}
}

type Lot struct {
	UserID       *string
	LotID        *string
	HoldingID    *string
	Shares       *float64
	CostPerShare *float64
	LotStatus    *uint32
	TradeDate    *uint64
	CreateTime   *uint64
	UpdateTime   *uint64
	Currency     *string
}

type LotOption func(l *Lot)

func WithLotID(lotID *string) LotOption {
	return func(l *Lot) {
		l.SetLotID(lotID)
	}
}

func WithLotShares(shares *float64) LotOption {
	return func(l *Lot) {
		l.SetShares(shares)
	}
}

func WithLotCostPerShare(costPerShare *float64) LotOption {
	return func(l *Lot) {
		l.SetCostPerShare(costPerShare)
	}
}

func WithLotStatus(lotStatus *uint32) LotOption {
	return func(l *Lot) {
		l.SetLotStatus(lotStatus)
	}
}

func WithLotTradeDate(tradeDate *uint64) LotOption {
	return func(l *Lot) {
		l.SetTradeDate(tradeDate)
	}
}

func WithLotCreateTime(createTime *uint64) LotOption {
	return func(l *Lot) {
		l.SetCreateTime(createTime)
	}
}

func WithLotUpdateTime(updateTime *uint64) LotOption {
	return func(l *Lot) {
		l.SetUpdateTime(updateTime)
	}
}

func WithLotCurrency(currency *string) LotOption {
	return func(l *Lot) {
		l.SetCurrency(currency)
	}
}

func NewLot(userID, holdingID string, opts ...LotOption) *Lot {
	now := uint64(time.Now().UnixMilli())
	l := &Lot{
		UserID:       goutil.String(userID),
		HoldingID:    goutil.String(holdingID),
		Shares:       goutil.Float64(0),
		CostPerShare: goutil.Float64(0),
		LotStatus:    goutil.Uint32(uint32(LotStatusNormal)),
		TradeDate:    goutil.Uint64(now),
		CreateTime:   goutil.Uint64(now),
		UpdateTime:   goutil.Uint64(now),
		Currency:     goutil.String(string(CurrencyUSD)),
	}

	for _, opt := range opts {
		opt(l)
	}

	l.validate()

	return l
}

func (l *Lot) validate() {}

func (l *Lot) Clone() *Lot {
	return NewLot(
		l.GetUserID(),
		l.GetHoldingID(),
		WithLotID(goutil.String(l.GetLotID())),
		WithLotShares(l.Shares),
		WithLotStatus(l.LotStatus),
		WithLotCostPerShare(l.CostPerShare),
		WithLotTradeDate(l.TradeDate),
		WithLotCreateTime(l.CreateTime),
		WithLotUpdateTime(l.UpdateTime),
		WithLotCurrency(l.Currency),
	)
}

type LotUpdate struct {
	Shares       *float64
	CostPerShare *float64
	TradeDate    *uint64
	LotStatus    *uint32
	UpdateTime   *uint64
}

func (l *Lot) ToLotUpdate(old *Lot) *LotUpdate {
	var (
		hasUpdate bool

		lu = &LotUpdate{
			UpdateTime: l.UpdateTime,
		}
	)

	if old.GetShares() != l.GetShares() {
		hasUpdate = true
		lu.Shares = l.Shares
	}

	if old.GetCostPerShare() != l.GetCostPerShare() {
		hasUpdate = true
		lu.CostPerShare = l.CostPerShare
	}

	if old.GetTradeDate() != l.GetTradeDate() {
		hasUpdate = true
		lu.TradeDate = l.TradeDate
	}

	if old.GetLotStatus() != l.GetLotStatus() {
		hasUpdate = true
		lu.LotStatus = l.LotStatus
	}

	if hasUpdate {
		return lu
	}

	return nil
}

func (l *Lot) Update(lus ...LotUpdateOption) *LotUpdate {
	if len(lus) == 0 {
		return nil
	}

	old := l.Clone()

	for _, lu := range lus {
		lu(l)
	}

	now := goutil.Uint64(uint64(time.Now().UnixMilli()))
	l.SetUpdateTime(now)

	return l.ToLotUpdate(old)
}

func (l *Lot) GetLotID() string {
	if l != nil && l.LotID != nil {
		return *l.LotID
	}
	return ""
}

func (l *Lot) SetLotID(lotID *string) {
	l.LotID = lotID
}

func (l *Lot) GetUserID() string {
	if l != nil && l.UserID != nil {
		return *l.UserID
	}
	return ""
}

func (l *Lot) SetUserID(userID *string) {
	l.UserID = userID
}

func (l *Lot) GetHoldingID() string {
	if l != nil && l.HoldingID != nil {
		return *l.HoldingID
	}
	return ""
}

func (l *Lot) SetHoldingID(holdingID *string) {
	l.HoldingID = holdingID
}

func (l *Lot) GetShares() float64 {
	if l != nil && l.Shares != nil {
		return *l.Shares
	}
	return 0
}

func (l *Lot) SetShares(shares *float64) {
	if shares != nil {
		s := util.RoundFloatToPreciseDP(*shares)
		l.Shares = goutil.Float64(s)
	}
}

func (l *Lot) GetCostPerShare() float64 {
	if l != nil && l.CostPerShare != nil {
		return *l.CostPerShare
	}
	return 0
}

func (l *Lot) SetCostPerShare(costPerShare *float64) {
	if costPerShare != nil {
		cps := util.RoundFloatToPreciseDP(*costPerShare)
		l.CostPerShare = goutil.Float64(cps)
	}
}

func (l *Lot) GetTradeDate() uint64 {
	if l != nil && l.TradeDate != nil {
		return *l.TradeDate
	}
	return 0
}

func (l *Lot) SetTradeDate(tradeDate *uint64) {
	l.TradeDate = tradeDate
}

func (l *Lot) GetLotStatus() uint32 {
	if l != nil && l.LotStatus != nil {
		return *l.LotStatus
	}
	return 0
}

func (l *Lot) SetLotStatus(lotStatus *uint32) {
	l.LotStatus = lotStatus
}

func (l *Lot) GetCreateTime() uint64 {
	if l != nil && l.CreateTime != nil {
		return *l.CreateTime
	}
	return 0
}

func (l *Lot) SetCreateTime(createTime *uint64) {
	l.CreateTime = createTime
}

func (l *Lot) GetUpdateTime() uint64 {
	if l != nil && l.UpdateTime != nil {
		return *l.UpdateTime
	}
	return 0
}

func (l *Lot) SetUpdateTime(updateTime *uint64) {
	l.UpdateTime = updateTime
}

func (l *Lot) GetCurrency() string {
	if l != nil && l.Currency != nil {
		return *l.Currency
	}
	return ""
}

func (l *Lot) SetCurrency(currency *string) {
	l.Currency = currency
}
