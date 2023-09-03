package transaction

import (
	"context"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/filter"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/common"
)

type UseCase interface {
	GetTransaction(ctx context.Context, req *GetTransactionRequest) (*GetTransactionResponse, error)
	GetTransactions(ctx context.Context, req *GetTransactionsRequest) (*GetTransactionsResponse, error)

	CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*CreateTransactionResponse, error)
	UpdateTransaction(ctx context.Context, req *UpdateTransactionRequest) (*UpdateTransactionResponse, error)
	DeleteTransaction(ctx context.Context, req *DeleteTransactionRequest) (*DeleteTransactionResponse, error)

	AggrTransactions(ctx context.Context, req *AggrTransactionsRequest) (*AggrTransactionsResponse, error)
}

type GetTransactionRequest struct {
	UserID        *string
	TransactionID *string
}

func (m *GetTransactionRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *GetTransactionRequest) GetTransactionID() string {
	if m != nil && m.TransactionID != nil {
		return *m.TransactionID
	}
	return ""
}

func (m *GetTransactionRequest) ToTransactionFilter() *repo.TransactionFilter {
	return repo.NewTransactionFilter(
		m.GetUserID(),
		repo.WithTransactionID(m.TransactionID),
	)
}

func (m *GetTransactionRequest) ToCategoryFilter(categoryID string) *repo.CategoryFilter {
	return &repo.CategoryFilter{
		UserID:         m.UserID,
		CategoryID:     goutil.String(categoryID),
		CategoryStatus: goutil.Uint32(uint32(entity.CategoryStatusNormal)),
	}
}

func (m *GetTransactionRequest) ToAccountFilter(accountID string) *repo.AccountFilter {
	return repo.NewAccountFilter(
		m.GetUserID(),
		repo.WithAccountID(goutil.String(accountID)),
		repo.WithAccountStatus(nil), // any status is ok
	)
}

type GetTransactionResponse struct {
	Transaction *entity.Transaction
}

func (m *GetTransactionResponse) GetTransaction() *entity.Transaction {
	if m != nil && m.Transaction != nil {
		return m.Transaction
	}
	return nil
}

type CreateTransactionRequest struct {
	UserID          *string
	AccountID       *string
	CategoryID      *string
	Amount          *float64
	Note            *string
	TransactionType *uint32
	TransactionTime *uint64
}

func (m *CreateTransactionRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *CreateTransactionRequest) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *CreateTransactionRequest) GetAccountID() string {
	if m != nil && m.AccountID != nil {
		return *m.AccountID
	}
	return ""
}

func (m *CreateTransactionRequest) GetAmount() float64 {
	if m != nil && m.Amount != nil {
		return *m.Amount
	}
	return 0
}

func (m *CreateTransactionRequest) GetNote() string {
	if m != nil && m.Note != nil {
		return *m.Note
	}
	return ""
}

func (m *CreateTransactionRequest) GetTransactionType() uint32 {
	if m != nil && m.TransactionType != nil {
		return *m.TransactionType
	}
	return 0
}

func (m *CreateTransactionRequest) GetTransactionTime() uint64 {
	if m != nil && m.TransactionTime != nil {
		return *m.TransactionTime
	}
	return 0
}

func (m *CreateTransactionRequest) ToTransactionEntity() *entity.Transaction {
	return entity.NewTransaction(
		m.GetUserID(),
		m.GetAccountID(),
		m.GetCategoryID(),
		entity.WithTransactionAmount(m.Amount),
		entity.WithTransactionNote(m.Note),
		entity.WithTransactionType(m.TransactionType),
		entity.WithTransactionTime(m.TransactionTime),
	)
}

func (m *CreateTransactionRequest) ToCategoryFilter() *repo.CategoryFilter {
	return &repo.CategoryFilter{
		UserID:         m.UserID,
		CategoryID:     m.CategoryID,
		CategoryStatus: goutil.Uint32(uint32(entity.CategoryStatusNormal)),
	}
}

func (m *CreateTransactionRequest) ToAccountFilter() *repo.AccountFilter {
	return repo.NewAccountFilter(
		m.GetUserID(),
		repo.WithAccountID(m.AccountID),
	)
}

type CreateTransactionResponse struct {
	Transaction *entity.Transaction
}

func (m *CreateTransactionResponse) GetTransaction() *entity.Transaction {
	if m != nil && m.Transaction != nil {
		return m.Transaction
	}
	return nil
}

type GetTransactionsRequest struct {
	UserID          *string
	AccountID       *string
	CategoryID      *string
	CategoryIDs     []string
	TransactionType *uint32
	TransactionTime *common.RangeFilter
	Paging          *common.Paging
}

func (m *GetTransactionsRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *GetTransactionsRequest) GetAccountID() string {
	if m != nil && m.AccountID != nil {
		return *m.AccountID
	}
	return ""
}

func (m *GetTransactionsRequest) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *GetTransactionsRequest) GetCategoryIDs() []string {
	if m != nil && m.CategoryIDs != nil {
		return m.CategoryIDs
	}
	return nil
}

func (m *GetTransactionsRequest) GetTransactionType() uint32 {
	if m != nil && m.TransactionType != nil {
		return *m.TransactionType
	}
	return 0
}

func (m *GetTransactionsRequest) GetTransactionTime() *common.RangeFilter {
	if m != nil && m.TransactionTime != nil {
		return m.TransactionTime
	}
	return nil
}

func (m *GetTransactionsRequest) GetPaging() *common.Paging {
	if m != nil && m.Paging != nil {
		return m.Paging
	}
	return nil
}

func (m *GetTransactionsRequest) ToTransactionFilter() *repo.TransactionFilter {
	tt := m.TransactionTime
	if tt == nil {
		tt = new(common.RangeFilter)
	}

	paging := m.Paging
	if m.Paging == nil {
		paging = new(common.Paging)
	}

	return repo.NewTransactionFilter(
		m.GetUserID(),
		repo.WithTransactionAccountID(m.AccountID),
		repo.WithTransactionCategoryID(m.CategoryID),
		repo.WithTransactionCategoryIDs(m.CategoryIDs),
		repo.WithTransactionType(m.TransactionType),
		repo.WithTransactionTimeGte(tt.Gte),
		repo.WithTransactionTimeLte(tt.Lte),
		repo.WithTransactionPaging(&repo.Paging{
			Limit: paging.Limit,
			Page:  paging.Page,
			Sorts: []filter.Sort{
				&repo.Sort{
					Field: goutil.String("transaction_time"),
					Order: goutil.String(config.OrderDesc),
				},
				&repo.Sort{
					Field: goutil.String("create_time"),
					Order: goutil.String(config.OrderDesc),
				},
			},
		}),
	)
}

func (m *GetTransactionsRequest) ToCategoryFilter(categoryIDs []string, categoryStatus uint32) *repo.CategoryFilter {
	return &repo.CategoryFilter{
		UserID:         m.UserID,
		CategoryIDs:    categoryIDs,
		CategoryStatus: goutil.Uint32(categoryStatus),
	}
}

func (m *GetTransactionsRequest) ToAccountFilter(accountIDs []string) *repo.AccountFilter {
	return repo.NewAccountFilter(
		m.GetUserID(),
		repo.WithAccountIDs(accountIDs),
		repo.WithAccountStatus(nil), // any status is ok
	)
}

type GetTransactionsResponse struct {
	Transactions []*entity.Transaction
	Paging       *common.Paging
}

func (m *GetTransactionsResponse) GetTransactions() []*entity.Transaction {
	if m != nil && m.Transactions != nil {
		return m.Transactions
	}
	return nil
}

func (m *GetTransactionsResponse) GetPaging() *common.Paging {
	if m != nil && m.Paging != nil {
		return m.Paging
	}
	return nil
}

type UpdateTransactionRequest struct {
	UserID          *string
	TransactionID   *string
	AccountID       *string
	CategoryID      *string
	TransactionType *uint32
	Note            *string
	Amount          *float64
	TransactionTime *uint64
}

func (m *UpdateTransactionRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *UpdateTransactionRequest) GetTransactionID() string {
	if m != nil && m.TransactionID != nil {
		return *m.TransactionID
	}
	return ""
}

func (m *UpdateTransactionRequest) GetAccountID() string {
	if m != nil && m.AccountID != nil {
		return *m.AccountID
	}
	return ""
}

func (m *UpdateTransactionRequest) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *UpdateTransactionRequest) GetTransactionType() uint32 {
	if m != nil && m.TransactionType != nil {
		return *m.TransactionType
	}
	return 0
}

func (m *UpdateTransactionRequest) GetAmount() float64 {
	if m != nil && m.Amount != nil {
		return *m.Amount
	}
	return 0
}

func (m *UpdateTransactionRequest) GetNote() string {
	if m != nil && m.Note != nil {
		return *m.Note
	}
	return ""
}

func (m *UpdateTransactionRequest) GetTransactionTime() uint64 {
	if m != nil && m.TransactionTime != nil {
		return *m.TransactionTime
	}
	return 0
}

func (m *UpdateTransactionRequest) ToTransactionFilter() *repo.TransactionFilter {
	return repo.NewTransactionFilter(
		m.GetUserID(),
		repo.WithTransactionID(m.TransactionID),
	)
}

func (m *UpdateTransactionRequest) ToCategoryFilter() *repo.CategoryFilter {
	return &repo.CategoryFilter{
		UserID:         m.UserID,
		CategoryID:     m.CategoryID,
		CategoryStatus: goutil.Uint32(uint32(entity.CategoryStatusNormal)),
	}
}

func (m *UpdateTransactionRequest) ToAccountFilter(accountID string) *repo.AccountFilter {
	return repo.NewAccountFilter(
		m.GetUserID(),
		repo.WithAccountID(goutil.String(accountID)),
	)
}

func (m *UpdateTransactionRequest) ToTransactionUpdate() *entity.TransactionUpdate {
	return entity.NewTransactionUpdate(
		entity.WithUpdateTransactionAmount(m.Amount),
		entity.WithUpdateTransactionTime(m.TransactionTime),
		entity.WithUpdateTransactionNote(m.Note),
		entity.WithUpdateTransactionAccountID(m.AccountID),
		entity.WithUpdateTransactionCategoryID(m.CategoryID),
		entity.WithUpdateTransactionType(m.TransactionType),
	)
}

type UpdateTransactionResponse struct {
	Transaction *entity.Transaction
}

func (m *UpdateTransactionResponse) GetTransaction() *entity.Transaction {
	if m != nil && m.Transaction != nil {
		return m.Transaction
	}
	return nil
}

type AggrTransactionsRequest struct {
	UserID           *string
	TransactionTime  *common.RangeFilter
	CategoryIDs      []string
	TransactionTypes []uint32
}

func (m *AggrTransactionsRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *AggrTransactionsRequest) GetTransactionTime() *common.RangeFilter {
	if m != nil && m.TransactionTime != nil {
		return m.TransactionTime
	}
	return nil
}

func (m *AggrTransactionsRequest) GetCategoryIDs() []string {
	if m != nil && m.CategoryIDs != nil {
		return m.CategoryIDs
	}
	return nil
}

func (m *AggrTransactionsRequest) GetTransactionTypes() []uint32 {
	if m != nil && m.TransactionTypes != nil {
		return m.TransactionTypes
	}
	return nil
}

func (m *AggrTransactionsRequest) ToTransactionFilter(userID string) *repo.TransactionFilter {
	tt := m.TransactionTime
	if tt == nil {
		tt = new(common.RangeFilter)
	}

	return repo.NewTransactionFilter(
		m.GetUserID(),
		repo.WithTransactionCategoryIDs(m.CategoryIDs),
		repo.WithTransactionTypes(m.TransactionTypes),
		repo.WithTransactionTimeGte(tt.Gte),
		repo.WithTransactionTimeLte(tt.Lte),
	)
}

func (m *AggrTransactionsRequest) ToCategoryFilter() *repo.CategoryFilter {
	return &repo.CategoryFilter{
		CategoryIDs:    m.CategoryIDs,
		CategoryStatus: goutil.Uint32(uint32(entity.CategoryStatusNormal)),
	}
}

type Aggr struct {
	Sum *float64
}

type AggrTransactionsResponse struct {
	Results map[string]*Aggr
}

func (m *AggrTransactionsResponse) GetResults() map[string]*Aggr {
	if m != nil && m.Results != nil {
		return m.Results
	}
	return nil
}

type DeleteTransactionRequest struct {
	UserID        *string
	TransactionID *string
}

func (m *DeleteTransactionRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *DeleteTransactionRequest) GetTransactionID() string {
	if m != nil && m.TransactionID != nil {
		return *m.TransactionID
	}
	return ""
}

func (m *DeleteTransactionRequest) ToTransactionFilter() *repo.TransactionFilter {
	return repo.NewTransactionFilter(
		m.GetUserID(),
		repo.WithTransactionID(m.TransactionID),
	)
}

func (m *DeleteTransactionRequest) ToAccountFilter(t *entity.Transaction) *repo.AccountFilter {
	return repo.NewAccountFilter(
		m.GetUserID(),
		repo.WithAccountID(t.AccountID),
	)
}

type DeleteTransactionResponse struct{}
