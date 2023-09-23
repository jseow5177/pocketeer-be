package presenter

import (
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/lot"
	"github.com/jseow5177/pockteer-be/util"
)

type Lot struct {
	LotID        *string `json:"lot_id,omitempty"`
	HoldingID    *string `json:"holding_id,omitempty"`
	Shares       *string `json:"shares,omitempty"`
	CostPerShare *string `json:"cost_per_share,omitempty"`
	LotStatus    *uint32 `json:"lot_status,omitempty"`
	TradeDate    *uint64 `json:"trade_date,omitempty"`
	Currency     *string `json:"currency,omitempty"`
	CreateTime   *uint64 `json:"create_time,omitempty"`
	UpdateTime   *uint64 `json:"update_time,omitempty"`
}

func (l *Lot) GetLotID() string {
	if l != nil && l.LotID != nil {
		return *l.LotID
	}
	return ""
}

func (l *Lot) GetHoldingID() string {
	if l != nil && l.HoldingID != nil {
		return *l.HoldingID
	}
	return ""
}

func (l *Lot) GetShares() string {
	if l != nil && l.Shares != nil {
		return *l.Shares
	}
	return ""
}

func (l *Lot) GetCurrency() string {
	if l != nil && l.Currency != nil {
		return *l.Currency
	}
	return ""
}

func (l *Lot) GetCostPerShare() string {
	if l != nil && l.CostPerShare != nil {
		return *l.CostPerShare
	}
	return ""
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

type DeleteLotRequest struct {
	LotID *string `json:"lot_id,omitempty"`
}

func (m *DeleteLotRequest) GetLotID() string {
	if m != nil && m.LotID != nil {
		return *m.LotID
	}
	return ""
}

func (m *DeleteLotRequest) ToUseCaseReq(userID string) *lot.DeleteLotRequest {
	return &lot.DeleteLotRequest{
		LotID:  m.LotID,
		UserID: goutil.String(userID),
	}
}

type DeleteLotResponse struct{}

func (m *DeleteLotResponse) Set(useCaseRes *lot.DeleteLotResponse) {}

type CreateLotRequest struct {
	HoldingID    *string `json:"holding_id,omitempty"`
	Shares       *string `json:"shares,omitempty"`
	CostPerShare *string `json:"cost_per_share,omitempty"`
	TradeDate    *uint64 `json:"trade_date,omitempty"`
}

func (m *CreateLotRequest) GetHoldingID() string {
	if m != nil && m.HoldingID != nil {
		return *m.HoldingID
	}
	return ""
}

func (m *CreateLotRequest) GetShares() string {
	if m != nil && m.Shares != nil {
		return *m.Shares
	}
	return ""
}

func (m *CreateLotRequest) GetCostPerShare() string {
	if m != nil && m.CostPerShare != nil {
		return *m.CostPerShare
	}
	return ""
}

func (m *CreateLotRequest) GetTradeDate() uint64 {
	if m != nil && m.TradeDate != nil {
		return *m.TradeDate
	}
	return 0
}

func (m *CreateLotRequest) ToUseCaseReq(userID string) *lot.CreateLotRequest {
	var shares *float64
	if m.Shares != nil {
		s, _ := util.MonetaryStrToFloat(m.GetShares())
		shares = goutil.Float64(s)
	}

	var costPerShare *float64
	if m.CostPerShare != nil {
		cps, _ := util.MonetaryStrToFloat(m.GetCostPerShare())
		costPerShare = goutil.Float64(cps)
	}

	return &lot.CreateLotRequest{
		UserID:       goutil.String(userID),
		Shares:       shares,
		CostPerShare: costPerShare,
		HoldingID:    m.HoldingID,
		TradeDate:    m.TradeDate,
	}
}

type CreateLotResponse struct {
	Lot *Lot `json:"lot,omitempty"`
}

func (m *CreateLotResponse) GetLot() *Lot {
	if m != nil && m.Lot != nil {
		return m.Lot
	}
	return nil
}

func (m *CreateLotResponse) Set(useCaseRes *lot.CreateLotResponse) {
	m.Lot = toLot(useCaseRes.Lot)
}

type UpdateLotRequest struct {
	LotID        *string `json:"lot_id,omitempty"`
	Shares       *string `json:"shares,omitempty"`
	CostPerShare *string `json:"cost_per_share,omitempty"`
	TradeDate    *uint64 `json:"trade_date,omitempty"`
}

func (m *UpdateLotRequest) GetLotID() string {
	if m != nil && m.LotID != nil {
		return *m.LotID
	}
	return ""
}

func (m *UpdateLotRequest) GetShares() string {
	if m != nil && m.Shares != nil {
		return *m.Shares
	}
	return ""
}

func (m *UpdateLotRequest) GetCostPerShare() string {
	if m != nil && m.CostPerShare != nil {
		return *m.CostPerShare
	}
	return ""
}

func (m *UpdateLotRequest) GetTradeDate() uint64 {
	if m != nil && m.TradeDate != nil {
		return *m.TradeDate
	}
	return 0
}

func (m *UpdateLotRequest) ToUseCaseReq(userID string) *lot.UpdateLotRequest {
	var shares *float64
	if m.Shares != nil {
		s, _ := util.MonetaryStrToFloat(m.GetShares())
		shares = goutil.Float64(s)
	}

	var costPerShare *float64
	if m.CostPerShare != nil {
		cps, _ := util.MonetaryStrToFloat(m.GetCostPerShare())
		costPerShare = goutil.Float64(cps)
	}

	return &lot.UpdateLotRequest{
		UserID:       goutil.String(userID),
		LotID:        m.LotID,
		Shares:       shares,
		CostPerShare: costPerShare,
	}
}

type UpdateLotResponse struct {
	Lot *Lot `json:"lot,omitempty"`
}

func (m *UpdateLotResponse) GetLot() *Lot {
	if m != nil && m.Lot != nil {
		return m.Lot
	}
	return nil
}

func (m *UpdateLotResponse) Set(useCaseRes *lot.UpdateLotResponse) {
	m.Lot = toLot(useCaseRes.Lot)
}

type GetLotRequest struct {
	UserID *string `json:"user_id,omitempty"`
	LotID  *string `json:"lot_id,omitempty"`
}

func (m *GetLotRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *GetLotRequest) GetLotID() string {
	if m != nil && m.LotID != nil {
		return *m.LotID
	}
	return ""
}

func (m *GetLotRequest) ToUseCaseReq(userID string) *lot.GetLotRequest {
	return &lot.GetLotRequest{
		UserID: goutil.String(userID),
		LotID:  m.LotID,
	}
}

type GetLotResponse struct {
	Lot *Lot `json:"lot,omitempty"`
}

func (m *GetLotResponse) GetLot() *Lot {
	if m != nil && m.Lot != nil {
		return m.Lot
	}
	return nil
}

func (m *GetLotResponse) Set(useCaseRes *lot.GetLotResponse) {
	m.Lot = toLot(useCaseRes.Lot)
}

type GetLotsRequest struct {
	HoldingID *string `json:"holding_id,omitempty"`
}

func (m *GetLotsRequest) GetHoldingID() string {
	if m != nil && m.HoldingID != nil {
		return *m.HoldingID
	}
	return ""
}

func (m *GetLotsRequest) ToUseCaseReq(userID string) *lot.GetLotsRequest {
	return &lot.GetLotsRequest{
		UserID:    goutil.String(userID),
		HoldingID: m.HoldingID,
	}
}

type GetLotsResponse struct {
	Lots []*Lot `json:"lots,omitempty"`
}

func (m *GetLotsResponse) GetLots() []*Lot {
	if m != nil && m.Lots != nil {
		return m.Lots
	}
	return nil
}

func (m *GetLotsResponse) Set(useCaseRes *lot.GetLotsResponse) {
	ls := make([]*Lot, 0)
	for _, l := range useCaseRes.Lots {
		ls = append(ls, toLot(l))
	}
	m.Lots = ls
}
