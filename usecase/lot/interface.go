package lot

import (
	"context"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/filter"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

type UseCase interface {
	GetLot(ctx context.Context, req *GetLotRequest) (*GetLotResponse, error)
	GetLots(ctx context.Context, req *GetLotsRequest) (*GetLotsResponse, error)

	CreateLot(ctx context.Context, req *CreateLotRequest) (*CreateLotResponse, error)
	UpdateLot(ctx context.Context, req *UpdateLotRequest) (*UpdateLotResponse, error)
	DeleteLot(ctx context.Context, req *DeleteLotRequest) (*DeleteLotResponse, error)
}

type DeleteLotRequest struct {
	UserID *string
	LotID  *string
}

func (m *DeleteLotRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *DeleteLotRequest) GetLotID() string {
	if m != nil && m.LotID != nil {
		return *m.LotID
	}
	return ""
}

func (m *DeleteLotRequest) ToLotFilter() *repo.LotFilter {
	return &repo.LotFilter{
		UserID:    m.UserID,
		LotID:     m.LotID,
		LotStatus: goutil.Uint32(uint32(entity.LotStatusNormal)),
	}
}

func (m *DeleteLotRequest) ToLotUpdate() *entity.LotUpdate {
	return entity.NewLotUpdate(
		entity.WithUpdateLotStatus(goutil.Uint32(uint32(entity.LotStatusDeleted))),
	)
}

type DeleteLotResponse struct{}

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
		UserID:    m.UserID,
		LotID:     m.LotID,
		LotStatus: goutil.Uint32(uint32(entity.LotStatusNormal)),
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
		LotStatus: goutil.Uint32(uint32(entity.LotStatusNormal)),
		Paging: &repo.Paging{
			Sorts: []filter.Sort{
				&repo.Sort{
					Field: goutil.String("trade_date"),
					Order: goutil.String(config.OrderDesc),
				},
				&repo.Sort{
					Field: goutil.String("create_time"),
					Order: goutil.String(config.OrderDesc),
				},
			},
		},
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

type UpdateLotRequest struct {
	UserID       *string
	LotID        *string
	Shares       *float64
	CostPerShare *float64
	TradeDate    *uint64
}

func (m *UpdateLotRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *UpdateLotRequest) GetLotID() string {
	if m != nil && m.LotID != nil {
		return *m.LotID
	}
	return ""
}

func (m *UpdateLotRequest) GetShares() float64 {
	if m != nil && m.Shares != nil {
		return *m.Shares
	}
	return 0
}

func (m *UpdateLotRequest) GetCostPerShare() float64 {
	if m != nil && m.CostPerShare != nil {
		return *m.CostPerShare
	}
	return 0
}

func (m *UpdateLotRequest) GetTradeDate() uint64 {
	if m != nil && m.TradeDate != nil {
		return *m.TradeDate
	}
	return 0
}

func (m *UpdateLotRequest) ToLotUpdate() *entity.LotUpdate {
	return entity.NewLotUpdate(
		entity.WithUpdateLotCostPerShare(m.CostPerShare),
		entity.WithUpdateLotShares(m.Shares),
		entity.WithUpdateLotTradeDate(m.TradeDate),
	)
}

func (m *UpdateLotRequest) ToLotFilter() *repo.LotFilter {
	return &repo.LotFilter{
		UserID:    m.UserID,
		LotID:     m.LotID,
		LotStatus: goutil.Uint32(uint32(entity.LotStatusNormal)),
	}
}

type UpdateLotResponse struct {
	Lot *entity.Lot
}

func (m *UpdateLotResponse) GetLot() *entity.Lot {
	if m != nil && m.Lot != nil {
		return m.Lot
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
		entity.WithLotShares(m.Shares),
		entity.WithLotCostPerShare(m.CostPerShare),
		entity.WithLotTradeDate(m.TradeDate),
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
