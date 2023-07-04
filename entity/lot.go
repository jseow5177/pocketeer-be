package entity

import (
	"time"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

type LotStatus uint32

const (
	LotStatusInvalid AccountStatus = iota
	LotStatusNormal
	LotStatusDeleted
)

type LotUpdate struct {
	Shares       *float64
	CostPerShare *float64
	TradeDate    *uint64
	UpdateTime   *uint64
}

func (lu *LotUpdate) GetShares() float64 {
	if lu != nil && lu.Shares != nil {
		return *lu.Shares
	}
	return 0
}

func (lu *LotUpdate) GetCostPerShare() float64 {
	if lu != nil && lu.CostPerShare != nil {
		return *lu.CostPerShare
	}
	return 0
}

func (lu *LotUpdate) GetTradeDate() uint64 {
	if lu != nil && lu.TradeDate != nil {
		return *lu.TradeDate
	}
	return 0
}

func (lu *LotUpdate) GetUpdateTime() uint64 {
	if lu != nil && lu.UpdateTime != nil {
		return *lu.UpdateTime
	}
	return 0
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
}

type LotOption = func(l *Lot)

func WithLotID(lotID *string) LotOption {
	return func(l *Lot) {
		l.LotID = lotID
	}
}

func WithShares(shares *float64) LotOption {
	return func(l *Lot) {
		l.Shares = shares
	}
}

func WithCostPerShare(costPerShare *float64) LotOption {
	return func(l *Lot) {
		l.CostPerShare = costPerShare
	}
}

func WithLotStatus(lotStatus *uint32) LotOption {
	return func(l *Lot) {
		l.LotStatus = lotStatus
	}
}

func WithTradeDate(tradeDate *uint64) LotOption {
	return func(l *Lot) {
		l.TradeDate = tradeDate
	}
}

func WithLotCreateTime(createTime *uint64) LotOption {
	return func(l *Lot) {
		l.CreateTime = createTime
	}
}

func WithLotUpdateTime(updateTime *uint64) LotOption {
	return func(l *Lot) {
		l.UpdateTime = updateTime
	}
}

func NewLot(userID, holdingID string, opts ...LotOption) *Lot {
	now := uint64(time.Now().Unix())
	l := &Lot{
		UserID:       goutil.String(userID),
		HoldingID:    goutil.String(holdingID),
		Shares:       goutil.Float64(0),
		CostPerShare: goutil.Float64(0),
		LotStatus:    goutil.Uint32(uint32(LotStatusNormal)),
		TradeDate:    goutil.Uint64(now),
		CreateTime:   goutil.Uint64(now),
		UpdateTime:   goutil.Uint64(now),
	}
	for _, opt := range opts {
		opt(l)
	}
	l.checkOpts()
	return l
}

func setLot(l *Lot, opts ...LotOption) {
	if l == nil {
		return
	}

	for _, opt := range opts {
		opt(l)
	}
}

func (l *Lot) checkOpts() {}

func (l *Lot) Update(lu *LotUpdate) (lotUpdate *LotUpdate, hasUpdate bool) {
	lotUpdate = new(LotUpdate)

	if lu.Shares != nil && lu.GetShares() != l.GetShares() {
		hasUpdate = true
		setLot(l, WithShares(lu.Shares))
	}

	if lu.CostPerShare != nil && lu.GetCostPerShare() != l.GetCostPerShare() {
		hasUpdate = true
		setLot(l, WithCostPerShare(lu.CostPerShare))
	}

	if lu.TradeDate != nil && lu.GetTradeDate() != l.GetTradeDate() {
		hasUpdate = true
		setLot(l, WithTradeDate(lu.TradeDate))
	}

	if hasUpdate {
		now := goutil.Uint64(uint64(time.Now().Unix()))
		setLot(l, WithLotUpdateTime(now))

		// check
		l.checkOpts()

		lotUpdate.UpdateTime = now

		if lu.Shares != nil {
			lotUpdate.Shares = l.Shares
		}

		if lu.CostPerShare != nil {
			lotUpdate.CostPerShare = l.CostPerShare
		}

		if lu.TradeDate != nil {
			lotUpdate.TradeDate = l.TradeDate
		}
	}

	return
}

func (l *Lot) GetLotID() string {
	if l != nil && l.LotID != nil {
		return *l.LotID
	}
	return ""
}

func (l *Lot) SetLotID(lotID string) {
	setLot(l, WithLotID(goutil.String(lotID)))
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
