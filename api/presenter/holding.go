package presenter

import (
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/holding"
	"github.com/jseow5177/pockteer-be/usecase/lot"
	"github.com/jseow5177/pockteer-be/util"
)

type Holding struct {
	HoldingID       *string `json:"holding_id,omitempty"`
	AccountID       *string `json:"account_id,omitempty"`
	Symbol          *string `json:"symbol,omitempty"`
	HoldingStatus   *uint32 `json:"holding_status,omitempty"`
	HoldingType     *uint32 `json:"holding_type,omitempty"`
	CreateTime      *uint64 `json:"create_time,omitempty"`
	UpdateTime      *uint64 `json:"update_time,omitempty"`
	TotalShares     *string `json:"total_shares,omitempty"`
	TotalCost       *string `json:"total_cost,omitempty"`
	AvgCostPerShare *string `json:"avg_cost_per_share,omitempty"`
	LatestValue     *string `json:"latest_value,omitempty"`
	Quote           *Quote  `json:"quote,omitempty"`
	Lots            []*Lot  `json:"lots,omitempty"`
	Currency        *string `json:"currency,omitempty"`
	Gain            *string `json:"gain,omitempty"`
	PercentGain     *string `json:"percent_gain,omitempty"`
}

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

func (h *Holding) GetTotalShares() string {
	if h != nil && h.TotalShares != nil {
		return *h.TotalShares
	}
	return ""
}

func (h *Holding) GetTotalCost() string {
	if h != nil && h.TotalCost != nil {
		return *h.TotalCost
	}
	return ""
}

func (h *Holding) GetLatestValue() string {
	if h != nil && h.LatestValue != nil {
		return *h.LatestValue
	}
	return ""
}

func (h *Holding) GetAvgCostPerShare() string {
	if h != nil && h.AvgCostPerShare != nil {
		return *h.AvgCostPerShare
	}
	return ""
}

func (h *Holding) GetQuote() *Quote {
	if h != nil && h.Quote != nil {
		return h.Quote
	}
	return nil
}

func (h *Holding) GetLots() []*Lot {
	if h != nil && h.Lots != nil {
		return h.Lots
	}
	return nil
}

func (h *Holding) GetCurrency() string {
	if h != nil && h.Currency != nil {
		return *h.Currency
	}
	return ""
}

func (h *Holding) GetGain() string {
	if h != nil && h.Gain != nil {
		return *h.Gain
	}
	return ""
}

func (h *Holding) GetPercentGain() string {
	if h != nil && h.PercentGain != nil {
		return *h.PercentGain
	}
	return ""
}

type UpdateHoldingRequest struct {
	HoldingID   *string             `json:"holding_id,omitempty"`
	TotalCost   *string             `json:"total_cost,omitempty"`
	LatestValue *string             `json:"latest_value,omitempty"`
	Symbol      *string             `json:"symbol,omitempty"`
	Currency    *string             `json:"currency,omitempty"`
	Lots        []*UpdateLotRequest `json:"lots,omitempty"`
}

func (m *UpdateHoldingRequest) GetHoldingID() string {
	if m != nil && m.HoldingID != nil {
		return *m.HoldingID
	}
	return ""
}

func (m *UpdateHoldingRequest) GetTotalCost() string {
	if m != nil && m.TotalCost != nil {
		return *m.TotalCost
	}
	return ""
}

func (m *UpdateHoldingRequest) GetLatestValue() string {
	if m != nil && m.LatestValue != nil {
		return *m.LatestValue
	}
	return ""
}

func (m *UpdateHoldingRequest) GetSymbol() string {
	if m != nil && m.Symbol != nil {
		return *m.Symbol
	}
	return ""
}

func (m *UpdateHoldingRequest) GetCurrency() string {
	if m != nil && m.Currency != nil {
		return *m.Currency
	}
	return ""
}

func (m *UpdateHoldingRequest) ToUseCaseReq(userID string) *holding.UpdateHoldingRequest {
	var totalCost *float64
	if m.TotalCost != nil {
		tc, _ := util.MonetaryStrToFloat(m.GetTotalCost())
		totalCost = goutil.Float64(tc)
	}

	var latestValue *float64
	if m.LatestValue != nil {
		lv, _ := util.MonetaryStrToFloat(m.GetLatestValue())
		latestValue = goutil.Float64(lv)
	}

	ls := make([]*lot.UpdateLotRequest, 0)
	for _, r := range m.Lots {
		ls = append(ls, r.ToUseCaseReq(userID))
	}

	return &holding.UpdateHoldingRequest{
		UserID:      goutil.String(userID),
		HoldingID:   m.HoldingID,
		TotalCost:   totalCost,
		LatestValue: latestValue,
		Symbol:      m.Symbol,
		Currency:    m.Currency,
		Lots:        ls,
	}
}

type UpdateHoldingResponse struct {
	Holding *Holding `json:"holding,omitempty"`
}

func (m *UpdateHoldingResponse) GetHolding() *Holding {
	if m != nil && m.Holding != nil {
		return m.Holding
	}
	return nil
}

func (m *UpdateHoldingResponse) Set(useCaseRes *holding.UpdateHoldingResponse) {
	m.Holding = toHolding(useCaseRes.Holding)
}

type CreateHoldingRequest struct {
	AccountID   *string             `json:"account_id,omitempty"`
	Symbol      *string             `json:"symbol,omitempty"`
	HoldingType *uint32             `json:"holding_type,omitempty"`
	TotalCost   *string             `json:"total_cost,omitempty"`
	LatestValue *string             `json:"latest_value,omitempty"`
	Currency    *string             `json:"currency,omitempty"`
	Lots        []*CreateLotRequest `json:"lots,omitempty"`
}

func (m *CreateHoldingRequest) GetAccountID() string {
	if m != nil && m.AccountID != nil {
		return *m.AccountID
	}
	return ""
}

func (m *CreateHoldingRequest) GetSymbol() string {
	if m != nil && m.Symbol != nil {
		return *m.Symbol
	}
	return ""
}

func (m *CreateHoldingRequest) GetCurrency() string {
	if m != nil && m.Currency != nil {
		return *m.Currency
	}
	return ""
}

func (m *CreateHoldingRequest) GetTotalCost() string {
	if m != nil && m.TotalCost != nil {
		return *m.TotalCost
	}
	return ""
}

func (m *CreateHoldingRequest) GetLatestValue() string {
	if m != nil && m.LatestValue != nil {
		return *m.LatestValue
	}
	return ""
}

func (m *CreateHoldingRequest) GetLots() []*CreateLotRequest {
	if m != nil && m.Lots != nil {
		return m.Lots
	}
	return nil
}

func (m *CreateHoldingRequest) ToUseCaseReq(userID string) *holding.CreateHoldingRequest {
	var totalCost *float64
	if m.TotalCost != nil {
		tc, _ := util.MonetaryStrToFloat(m.GetTotalCost())
		totalCost = goutil.Float64(tc)
	}

	var latestValue *float64
	if m.LatestValue != nil {
		lv, _ := util.MonetaryStrToFloat(m.GetLatestValue())
		latestValue = goutil.Float64(lv)
	}

	ls := make([]*lot.CreateLotRequest, 0)
	for _, r := range m.Lots {
		ls = append(ls, r.ToUseCaseReq(userID))
	}

	return &holding.CreateHoldingRequest{
		UserID:      goutil.String(userID),
		AccountID:   m.AccountID,
		Symbol:      m.Symbol,
		HoldingType: m.HoldingType,
		Currency:    m.Currency,
		TotalCost:   totalCost,
		LatestValue: latestValue,
		Lots:        ls,
	}
}

type CreateHoldingResponse struct {
	Holding *Holding `json:"holding,omitempty"`
}

func (m *CreateHoldingResponse) GetHolding() *Holding {
	if m != nil && m.Holding != nil {
		return m.Holding
	}
	return nil
}

func (m *CreateHoldingResponse) Set(useCaseRes *holding.CreateHoldingResponse) {
	m.Holding = toHolding(useCaseRes.Holding)
}

type GetHoldingRequest struct {
	HoldingID *string `json:"holding_id,omitempty"`
}

func (m *GetHoldingRequest) GetHoldingID() string {
	if m != nil && m.HoldingID != nil {
		return *m.HoldingID
	}
	return ""
}

func (m *GetHoldingRequest) ToUseCaseReq(userID string) *holding.GetHoldingRequest {
	return &holding.GetHoldingRequest{
		UserID:    goutil.String(userID),
		HoldingID: m.HoldingID,
	}
}

type GetHoldingResponse struct {
	Holding *Holding `json:"holding,omitempty"`
}

func (m *GetHoldingResponse) GetHolding() *Holding {
	if m != nil && m.Holding != nil {
		return m.Holding
	}
	return nil
}

func (m *GetHoldingResponse) Set(useCaseRes *holding.GetHoldingResponse) {
	m.Holding = toHolding(useCaseRes.Holding)
}

type DeleteHoldingRequest struct {
	HoldingID *string `json:"holding_id,omitempty"`
}

func (m *DeleteHoldingRequest) GetAccountID() string {
	if m != nil && m.HoldingID != nil {
		return *m.HoldingID
	}
	return ""
}

func (m *DeleteHoldingRequest) ToUseCaseReq(userID string) *holding.DeleteHoldingRequest {
	return &holding.DeleteHoldingRequest{
		UserID:    goutil.String(userID),
		HoldingID: m.HoldingID,
	}
}

type DeleteHoldingResponse struct{}

func (m *DeleteHoldingResponse) Set(useCaseRes *holding.DeleteHoldingResponse) {}
