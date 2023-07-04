package lot

import (
	"context"

	"github.com/jseow5177/pockteer-be/entity"
)

type UseCase interface {
	CreateLot(ctx context.Context, req *CreateLotRequest) (*CreateLotResponse, error)
}

type CreateLotRequest struct {
	UserID       *string
	HoldingID    *string
	Shares       *float64
	CostPerShare *float64
	LotStatus    *uint32
	TradeDate    *uint64
}

func (m *CreateLotRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *CreateLotRequest) GetHoldingID() string {
	if m != nil && m.HoldingID != nil {
		return *m.HoldingID
	}
	return ""
}

func (m *CreateLotRequest) GetShares() float64 {
	if m != nil && m.Shares != nil {
		return *m.Shares
	}
	return 0
}

func (m *CreateLotRequest) GetCostPerShare() float64 {
	if m != nil && m.CostPerShare != nil {
		return *m.CostPerShare
	}
	return 0
}

func (m *CreateLotRequest) GetLotStatus() uint32 {
	if m != nil && m.LotStatus != nil {
		return *m.LotStatus
	}
	return 0
}

func (m *CreateLotRequest) GetTradeDate() uint64 {
	if m != nil && m.TradeDate != nil {
		return *m.TradeDate
	}
	return 0
}

type CreateLotResponse struct {
	Lot *entity.Lot
}

func (m *CreateLotResponse) GetLot() *entity.Lot {
	if m != nil && m.Lot != nil {
		return m.Lot
	}
	return nil
}
