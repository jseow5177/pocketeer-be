package presenter

import (
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/holding"
)

type Holding struct {
	HoldingID     *string `json:"_id,omitempty"`
	AccountID     *string `json:"account_id,omitempty"`
	Symbol        *string `json:"symbol,omitempty"`
	HoldingStatus *uint32 `json:"holding_status,omitempty"`
	HoldingType   *uint32 `json:"holding_type,omitempty"`
	CreateTime    *uint64 `json:"create_time,omitempty"`
	UpdateTime    *uint64 `json:"update_time,omitempty"`
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

type CreateHoldingRequest struct {
	AccountID   *string `json:"account_id,omitempty"`
	Symbol      *string `json:"symbol,omitempty"`
	HoldingType *uint32 `json:"holding_type,omitempty"`
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

func (m *CreateHoldingRequest) GetHoldingType() uint32 {
	if m != nil && m.HoldingType != nil {
		return *m.HoldingType
	}
	return 0
}

func (m *CreateHoldingRequest) ToUseCaseReq(userID string) *holding.CreateHoldingRequest {
	return &holding.CreateHoldingRequest{
		UserID:      goutil.String(userID),
		AccountID:   m.AccountID,
		Symbol:      m.Symbol,
		HoldingType: m.HoldingType,
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
