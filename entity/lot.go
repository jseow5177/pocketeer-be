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

type LotUpdate struct {
	Shares       *float64
	CostPerShare *float64
	TradeDate    *uint64
	LotStatus    *uint32
	UpdateTime   *uint64
}

func (lu *LotUpdate) GetShares() float64 {
	if lu != nil && lu.Shares != nil {
		return *lu.Shares
	}
	return 0
}

func (lu *LotUpdate) SetShares(shares *float64) {
	if shares != nil {
		s := util.RoundFloatToPreciseDP(*shares)
		lu.Shares = goutil.Float64(s)
	}
}

func (lu *LotUpdate) GetCostPerShare() float64 {
	if lu != nil && lu.CostPerShare != nil {
		return *lu.CostPerShare
	}
	return 0
}

func (lu *LotUpdate) SetCostPerShare(costPerShare *float64) {
	lu.CostPerShare = costPerShare
}

func (lu *LotUpdate) GetLotStatus() uint32 {
	if lu != nil && lu.LotStatus != nil {
		return *lu.LotStatus
	}
	return 0
}

func (lu *LotUpdate) SetLotStatus(lotStatus *uint32) {
	lu.LotStatus = lotStatus
}

func (lu *LotUpdate) GetTradeDate() uint64 {
	if lu != nil && lu.TradeDate != nil {
		return *lu.TradeDate
	}
	return 0
}

func (lu *LotUpdate) SetTradeDate(tradeDate *uint64) {
	lu.TradeDate = tradeDate
}

func (lu *LotUpdate) GetUpdateTime() uint64 {
	if lu != nil && lu.UpdateTime != nil {
		return *lu.UpdateTime
	}
	return 0
}

func (lu *LotUpdate) SetUpdateTime(updateTime *uint64) {
	lu.UpdateTime = updateTime
}

func WithUpdateLotShares(shares *float64) LotUpdateOption {
	return func(lu *LotUpdate) {
		lu.SetShares(shares)
	}
}

func WithUpdateLotCostPerShare(costPerShare *float64) LotUpdateOption {
	return func(lu *LotUpdate) {
		lu.SetCostPerShare(costPerShare)
	}
}

func WithUpdateLotTradeDate(tradeDate *uint64) LotUpdateOption {
	return func(lu *LotUpdate) {
		lu.SetTradeDate(tradeDate)
	}
}

func WithUpdateLotStatus(lotStatus *uint32) LotUpdateOption {
	return func(lu *LotUpdate) {
		lu.SetLotStatus(lotStatus)
	}
}

func NewLotUpdate(opts ...LotUpdateOption) *LotUpdate {
	lu := new(LotUpdate)
	for _, opt := range opts {
		opt(lu)
	}
	return lu
}

type LotUpdateOption = func(lu *LotUpdate)

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
	}
	for _, opt := range opts {
		opt(l)
	}
	l.checkOpts()
	return l
}

func (l *Lot) checkOpts() {}

func (l *Lot) Update(lu *LotUpdate) (lotUpdate *LotUpdate, hasUpdate bool) {
	lotUpdate = new(LotUpdate)

	if lu.Shares != nil && lu.GetShares() != l.GetShares() {
		hasUpdate = true
		l.SetShares(lu.Shares)

		defer func() {
			lotUpdate.SetShares(l.Shares)
		}()
	}

	if lu.CostPerShare != nil && lu.GetCostPerShare() != l.GetCostPerShare() {
		hasUpdate = true
		l.SetCostPerShare(lu.CostPerShare)

		defer func() {
			lotUpdate.SetCostPerShare(l.CostPerShare)
		}()
	}

	if lu.TradeDate != nil && lu.GetTradeDate() != l.GetTradeDate() {
		hasUpdate = true
		l.SetTradeDate(lu.TradeDate)

		defer func() {
			lotUpdate.SetTradeDate(l.TradeDate)
		}()
	}

	if lu.LotStatus != nil && lu.GetLotStatus() != l.GetLotStatus() {
		hasUpdate = true
		l.SetLotStatus(lu.LotStatus)

		defer func() {
			lotUpdate.SetLotStatus(l.LotStatus)
		}()
	}

	if !hasUpdate {
		return
	}

	now := goutil.Uint64(uint64(time.Now().UnixMilli()))
	l.SetUpdateTime(now)

	l.checkOpts()

	lotUpdate.SetUpdateTime(now)

	return
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
	l.CostPerShare = costPerShare
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
