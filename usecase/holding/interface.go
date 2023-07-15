package holding

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
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
	AvgCost     *float64
	LatestValue *float64
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

func (m *UpdateHoldingRequest) GetAvgCost() float64 {
	if m != nil && m.AvgCost != nil {
		return *m.AvgCost
	}
	return 0
}

func (m *UpdateHoldingRequest) GetLatestValue() float64 {
	if m != nil && m.LatestValue != nil {
		return *m.LatestValue
	}
	return 0
}

func (m *UpdateHoldingRequest) ToHoldingFilter() *repo.HoldingFilter {
	return &repo.HoldingFilter{
		UserID:    m.UserID,
		HoldingID: m.HoldingID,
	}
}

func (m *UpdateHoldingRequest) ToHoldingUpdate() *entity.HoldingUpdate {
	return &entity.HoldingUpdate{
		AvgCost:     m.AvgCost,
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
	AvgCost     *float64
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

func (m *CreateHoldingRequest) GetAvgCost() float64 {
	if m != nil && m.AvgCost != nil {
		return *m.AvgCost
	}
	return 0
}

func (m *CreateHoldingRequest) GetLatestValue() float64 {
	if m != nil && m.LatestValue != nil {
		return *m.LatestValue
	}
	return 0
}

func (m *CreateHoldingRequest) ToAccountFilter(userID string) *repo.AccountFilter {
	return repo.NewAccountFilter(userID, repo.WitAccountID(m.AccountID))
}

func (m *CreateHoldingRequest) ToHoldingEntity() (*entity.Holding, error) {
	return entity.NewHolding(
		m.GetUserID(),
		m.GetAccountID(),
		m.GetSymbol(),
		entity.WithHoldingType(m.HoldingType),
		entity.WithHoldingAvgCost(m.AvgCost),
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
