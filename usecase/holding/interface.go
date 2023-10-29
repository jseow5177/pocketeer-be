package holding

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/lot"
)

type UseCase interface {
	GetHolding(ctx context.Context, req *GetHoldingRequest) (*GetHoldingResponse, error)

	CreateHolding(ctx context.Context, req *CreateHoldingRequest) (*CreateHoldingResponse, error)
	UpdateHolding(ctx context.Context, req *UpdateHoldingRequest) (*UpdateHoldingResponse, error)
	DeleteHolding(ctx context.Context, req *DeleteHoldingRequest) (*DeleteHoldingResponse, error)
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
	return repo.NewHoldingFilter(
		repo.WithHoldingUserID(m.UserID),
		repo.WithHoldingID(m.HoldingID),
	)
}

func (m *GetHoldingRequest) ToQuoteFilter(symbol string) *repo.QuoteFilter {
	return repo.NewQuoteFilter(
		repo.WithQuoteSymbol(goutil.String(symbol)),
	)
}

func (m *GetHoldingRequest) ToSecurityFilter(symbol string) *repo.SecurityFilter {
	return repo.NewSecurityFilter(
		repo.WithSecuritySymbol(goutil.String(symbol)),
	)
}

func (m *GetHoldingRequest) ToLotFilter() *repo.LotFilter {
	return repo.NewLotFilter(
		m.GetUserID(),
		repo.WithLotHoldingID(m.HoldingID),
	)
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
	Lots        []*lot.UpdateLotRequest
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
	return repo.NewHoldingFilter(
		repo.WithHoldingUserID(m.UserID),
		repo.WithHoldingID(m.HoldingID),
	)
}

func (m *UpdateHoldingRequest) ToLotFilters() *repo.LotFilter {
	lotIDs := make([]string, 0)
	for _, lot := range m.Lots {
		lotIDs = append(lotIDs, lot.GetLotID())
	}
	return repo.NewLotFilter(
		m.GetUserID(),
		repo.WithLotHoldingID(m.HoldingID),
		repo.WithLotIDs(lotIDs),
	)
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
	Currency    *string
	HoldingType *uint32
	TotalCost   *float64
	LatestValue *float64
	Lots        []*lot.CreateLotRequest
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

func (m *CreateHoldingRequest) GetCurrency() string {
	if m != nil && m.Currency != nil {
		return *m.Currency
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
	return repo.NewAccountFilter(
		m.GetUserID(),
		repo.WithAccountID(m.AccountID),
	)
}

func (m *CreateHoldingRequest) ToSecurityFilter() *repo.SecurityFilter {
	return repo.NewSecurityFilter(
		repo.WithSecuritySymbol(m.Symbol),
	)
}

func (m *CreateHoldingRequest) ToQuoteFilter() *repo.QuoteFilter {
	return repo.NewQuoteFilter(
		repo.WithQuoteSymbol(m.Symbol),
	)
}

func (m *CreateHoldingRequest) ToHoldingEntity(currency string) (*entity.Holding, error) {
	return entity.NewHolding(
		m.GetUserID(),
		m.GetAccountID(),
		m.GetSymbol(),
		entity.WithHoldingType(m.HoldingType),
		entity.WithHoldingTotalCost(m.TotalCost),
		entity.WithHoldingLatestValue(m.LatestValue),
		entity.WithHoldingCurrency(goutil.String(currency)),
		entity.WithHoldingLots(m.ToLotEntities(currency)),
	)
}

func (m *CreateHoldingRequest) ToLotEntities(currency string) []*entity.Lot {
	ls := make([]*entity.Lot, 0)
	for _, r := range m.Lots {
		ls = append(ls, r.ToLotEntity(currency))
	}
	return ls
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

type DeleteHoldingRequest struct {
	UserID    *string
	HoldingID *string
}

func (m *DeleteHoldingRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *DeleteHoldingRequest) GetAccountID() string {
	if m != nil && m.HoldingID != nil {
		return *m.HoldingID
	}
	return ""
}

func (m *DeleteHoldingRequest) ToHoldingFilter() *repo.HoldingFilter {
	return repo.NewHoldingFilter(
		repo.WithHoldingUserID(m.UserID),
		repo.WithHoldingID(m.HoldingID),
	)
}

func (m *DeleteHoldingRequest) ToLotFilter() *repo.LotFilter {
	return repo.NewLotFilter(
		m.GetUserID(),
		repo.WithLotHoldingID(m.HoldingID),
	)
}

type DeleteHoldingResponse struct{}
