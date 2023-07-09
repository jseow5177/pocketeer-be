package holding

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
)

type UseCase interface {
	GetHolding(ctx context.Context, req *GetHoldingRequest) (*GetHoldingResponse, error)
	CreateHolding(ctx context.Context, req *CreateHoldingRequest) (*CreateHoldingResponse, error)
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

func (m *GetHoldingRequest) ToLotFilter() *repo.LotFilter {
	return &repo.LotFilter{
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

type CreateHoldingRequest struct {
	UserID      *string
	AccountID   *string
	Symbol      *string
	HoldingType *uint32
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

func (m *CreateHoldingRequest) ToAccountFilter(userID string) *repo.AccountFilter {
	return repo.NewAccountFilter(userID, repo.WitAccountID(m.AccountID))
}

func (m *CreateHoldingRequest) ToHoldingEntity() *entity.Holding {
	return entity.NewHolding(
		m.GetUserID(),
		m.GetAccountID(),
		m.GetSymbol(),
		entity.WithHoldingType(m.HoldingType),
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
