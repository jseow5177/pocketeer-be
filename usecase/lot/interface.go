package lot

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

type UseCase interface {
	GetLot(ctx context.Context, req *GetLotRequest) (*GetLotResponse, error)
	GetLots(ctx context.Context, req *GetLotsRequest) (*GetLotsResponse, error)

	CreateLot(ctx context.Context, req *CreateLotRequest) (*CreateLotResponse, error)
}

type GetLotRequest struct {
	UserID *string
	LotID  *string
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

func (m *GetLotRequest) ToLotFilter() *repo.LotFilter {
	return &repo.LotFilter{
		UserID: m.UserID,
		LotID:  m.LotID,
	}
}

type GetLotResponse struct {
	Lot *entity.Lot
}

func (m *GetLotResponse) GetLot() *entity.Lot {
	if m != nil && m.Lot != nil {
		return m.Lot
	}
	return nil
}

type GetLotsRequest struct {
	UserID    *string
	HoldingID *string
}

func (m *GetLotsRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *GetLotsRequest) GetHoldingID() string {
	if m != nil && m.HoldingID != nil {
		return *m.HoldingID
	}
	return ""
}

func (m *GetLotsRequest) ToLotFilter() *repo.LotFilter {
	return &repo.LotFilter{
		UserID:    m.UserID,
		HoldingID: m.HoldingID,
	}
}

type GetLotsResponse struct {
	Lots []*entity.Lot
}

func (m *GetLotsResponse) GetLots() []*entity.Lot {
	if m != nil && m.Lots != nil {
		return m.Lots
	}
	return nil
}

type CreateLotRequest struct {
	UserID       *string
	HoldingID    *string
	Shares       *float64
	CostPerShare *float64
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

func (m *CreateLotRequest) GetTradeDate() uint64 {
	if m != nil && m.TradeDate != nil {
		return *m.TradeDate
	}
	return 0
}

func (m *CreateLotRequest) ToHoldingFilter() *repo.HoldingFilter {
	return &repo.HoldingFilter{
		UserID:    m.UserID,
		HoldingID: m.HoldingID,
	}
}

func (m *CreateLotRequest) ToLotEntity() *entity.Lot {
	return entity.NewLot(
		m.GetUserID(),
		m.GetHoldingID(),
		entity.WithShares(m.Shares),
		entity.WithCostPerShare(m.CostPerShare),
		entity.WithTradeDate(m.TradeDate),
		entity.WithLotStatus(goutil.Uint32(uint32(entity.LotStatusNormal))),
	)
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
