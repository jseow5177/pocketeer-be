package holding

import (
	"context"

	"github.com/jseow5177/pockteer-be/entity"
)

type UseCase interface {
	CreateHolding(ctx context.Context, req *CreateHoldingRequest) (*CreateHoldingResponse, error)
}

type CreateHoldingRequest struct {
	AccountID   *string
	Symbol      *string
	HoldingType *uint32
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

func (m *CreateHoldingRequest) ToHoldingEntity() *entity.Holding {
	return entity.NewHolding(
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
