package holding

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

type UseCase interface {
	GetHolding(ctx context.Context, req *GetHoldingRequest) (*GetHoldingResponse, error)
	GetHoldings(ctx context.Context, req *GetHoldingsRequest) (*GetHoldingsResponse, error)

	CreateHolding(ctx context.Context, req *CreateHoldingRequest) (*CreateHoldingResponse, error)
	UpdateHolding(ctx context.Context, req *UpdateHoldingRequest) (*UpdateHoldingResponse, error)
}

type GetHoldingsRequest struct {
	UserID    *string
	AccountID *string
}

func (m *GetHoldingsRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *GetHoldingsRequest) GetAccountID() string {
	if m != nil && m.AccountID != nil {
		return *m.AccountID
	}
	return ""
}

func (m *GetHoldingsRequest) ToHoldingFilter() *repo.HoldingFilter {
	return &repo.HoldingFilter{
		UserID:    m.UserID,
		AccountID: m.AccountID,
	}
}

func (m *GetHoldingsRequest) ToQuoteFilter(symbol string) *repo.QuoteFilter {
	return &repo.QuoteFilter{
		Symbol: goutil.String(symbol),
	}
}

type GetHoldingsResponse struct {
	Holdings []*entity.Holding
}

func (m *GetHoldingsResponse) GetHoldings() []*entity.Holding {
	if m != nil && m.Holdings != nil {
		return m.Holdings
	}
	return nil
}

type GetHoldingRequest struct {
	UserID    *string
	HoldingID *string
}

func (m *GetHoldingRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *GetHoldingRequest) GetHoldingID() string {
	if m != nil && m.HoldingID != nil {
		return *m.HoldingID
	}
	return ""
}

func (m *GetHoldingRequest) ToHoldingFilter() *repo.HoldingFilter {
	return &repo.HoldingFilter{
		UserID:    m.UserID,
		HoldingID: m.HoldingID,
	}
}

func (m *GetHoldingRequest) ToQuoteFilter(symbol string) *repo.QuoteFilter {
	return &repo.QuoteFilter{
		Symbol: goutil.String(symbol),
	}
}

type GetHoldingResponse struct {
	Holding *entity.Holding
}

func (m *GetHoldingResponse) GetHolding() *entity.Holding {
	if m != nil && m.Holding != nil {
		return m.Holding
	}
	return nil
}

type UpdateHoldingRequest struct {
	UserID      *string
	HoldingID   *string
	TotalCost   *float64
	LatestValue *float64
	Symbol      *string
}

func (m *UpdateHoldingRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *UpdateHoldingRequest) GetHoldingID() string {
	if m != nil && m.HoldingID != nil {
		return *m.HoldingID
	}
	return ""
}

func (m *UpdateHoldingRequest) GetTotalCost() float64 {
	if m != nil && m.TotalCost != nil {
		return *m.TotalCost
	}
	return 0
}

func (m *UpdateHoldingRequest) GetLatestValue() float64 {
	if m != nil && m.LatestValue != nil {
		return *m.LatestValue
	}
	return 0
}

func (m *UpdateHoldingRequest) GetSymbol() string {
	if m != nil && m.Symbol != nil {
		return *m.Symbol
	}
	return ""
}

func (m *UpdateHoldingRequest) ToHoldingFilter() *repo.HoldingFilter {
	return &repo.HoldingFilter{
		UserID:    m.UserID,
		HoldingID: m.HoldingID,
	}
}

func (m *UpdateHoldingRequest) ToHoldingUpdate() *entity.HoldingUpdate {
	return &entity.HoldingUpdate{
		TotalCost:   m.TotalCost,
		LatestValue: m.LatestValue,
	}
}

type UpdateHoldingResponse struct {
	Holding *entity.Holding
}

func (m *UpdateHoldingResponse) GetHolding() *entity.Holding {
	if m != nil && m.Holding != nil {
		return m.Holding
	}
	return nil
}

type CreateHoldingRequest struct {
	UserID      *string
	AccountID   *string
	Symbol      *string
	HoldingType *uint32
	TotalCost   *float64
	LatestValue *float64
}

func (m *CreateHoldingRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
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

func (m *CreateHoldingRequest) GetTotalCost() float64 {
	if m != nil && m.TotalCost != nil {
		return *m.TotalCost
	}
	return 0
}

func (m *CreateHoldingRequest) GetLatestValue() float64 {
	if m != nil && m.LatestValue != nil {
		return *m.LatestValue
	}
	return 0
}

func (m *CreateHoldingRequest) ToAccountFilter() *repo.AccountFilter {
	return repo.NewAccountFilter(m.GetUserID(), repo.WithAccountID(m.AccountID))
}

func (m *CreateHoldingRequest) ToSecurityFilter() *repo.SecurityFilter {
	return &repo.SecurityFilter{
		Symbol: m.Symbol,
	}
}

func (m *CreateHoldingRequest) ToHoldingEntity() (*entity.Holding, error) {
	return entity.NewHolding(
		m.GetUserID(),
		m.GetAccountID(),
		m.GetSymbol(),
		entity.WithHoldingType(m.HoldingType),
		entity.WithHoldingTotalCost(m.TotalCost),
		entity.WithHoldingLatestValue(m.LatestValue),
	)
}

type CreateHoldingResponse struct {
	Holding *entity.Holding
}

func (m *CreateHoldingResponse) GetHolding() *entity.Holding {
	if m != nil && m.Holding != nil {
		return m.Holding
	}
	return nil
}
