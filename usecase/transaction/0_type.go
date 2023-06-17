package transaction

import (
	"context"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/filter"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/account"
	"github.com/jseow5177/pockteer-be/usecase/category"
	"github.com/jseow5177/pockteer-be/usecase/common"
)

type UseCase interface {
	GetTransaction(ctx context.Context, req *GetTransactionRequest) (*GetTransactionResponse, error)
	GetTransactions(ctx context.Context, req *GetTransactionsRequest) (*GetTransactionsResponse, error)

	CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*CreateTransactionResponse, error)
	UpdateTransaction(ctx context.Context, req *UpdateTransactionRequest) (*UpdateTransactionResponse, error)

	AggrTransactions(ctx context.Context, req *AggrTransactionsRequest) (*AggrTransactionsResponse, error)
}

type Transaction struct {
	*entity.Transaction

	Category *entity.Category
	Account  *entity.Account
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
	return &repo.TransactionFilter{
		UserID:        m.UserID,
		TransactionID: m.TransactionID,
	}
}

func (m *GetTransactionRequest) ToGetCategoryRequest(categoryID string) *category.GetCategoryRequest {
	return &category.GetCategoryRequest{
		UserID:     m.UserID,
		CategoryID: goutil.String(categoryID),
	}
}

func (m *GetTransactionRequest) ToGetAccountRequest(accountID string) *account.GetAccountRequest {
	return &account.GetAccountRequest{
		UserID:    m.UserID,
		AccountID: goutil.String(accountID),
	}
}

type GetTransactionResponse struct {
	Transaction *Transaction
}

func (m *GetTransactionResponse) GetTransaction() *Transaction {
	if m != nil && m.Transaction != nil {
		return m.Transaction
	}
	return nil
}

type CreateTransactionRequest struct {
	UserID          *string
	AccountID       *string
	CategoryID      *string
	Amount          *string
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

func (m *CreateTransactionRequest) GetAmount() string {
	if m != nil && m.Amount != nil {
		return *m.Amount
	}
	return ""
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
	t := &entity.Transaction{
		UserID:          m.UserID,
		CategoryID:      m.CategoryID,
		Note:            m.Note,
		TransactionType: m.TransactionType,
		TransactionTime: m.TransactionTime,
	}

	t.SetAmount(m.GetAmount())

	return t
}

func (m *CreateTransactionRequest) ToGetCategoryRequest() *category.GetCategoryRequest {
	return &category.GetCategoryRequest{
		UserID:     m.UserID,
		CategoryID: m.CategoryID,
	}
}

type CreateTransactionResponse struct {
	Transaction *Transaction
}

func (m *CreateTransactionResponse) GetTransaction() *Transaction {
	if m != nil && m.Transaction != nil {
		return m.Transaction
	}
	return nil
}

type GetTransactionsRequest struct {
	UserID          *string
	CategoryID      *string
	TransactionType *uint32
	TransactionTime *common.UInt64Filter
	Paging          *common.Paging
}

func (m *GetTransactionsRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *GetTransactionsRequest) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *GetTransactionsRequest) GetTransactionType() uint32 {
	if m != nil && m.TransactionType != nil {
		return *m.TransactionType
	}
	return 0
}

func (m *GetTransactionsRequest) GetTransactionTime() *common.UInt64Filter {
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

func (m *GetTransactionsRequest) ToGetCategoryRequest() *category.GetCategoryRequest {
	return &category.GetCategoryRequest{
		CategoryID: m.CategoryID,
	}
}

func (m *GetTransactionsRequest) ToTransactionFilter() *repo.TransactionFilter {
	tt := m.TransactionTime
	if tt == nil {
		tt = new(common.UInt64Filter)
	}

	paging := m.Paging
	if m.Paging == nil {
		paging = new(common.Paging)
	}

	return &repo.TransactionFilter{
		UserID:             m.UserID,
		CategoryID:         m.CategoryID,
		TransactionType:    m.TransactionType,
		TransactionTimeGte: tt.Gte,
		TransactionTimeLte: tt.Lte,
		Paging: &repo.Paging{
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
		},
	}
}

type GetTransactionsResponse struct {
	Transactions []*Transaction
	Paging       *common.Paging
}

func (m *GetTransactionsResponse) GetTransactions() []*Transaction {
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
	CategoryID      *string
	Note            *string
	Amount          *string
	TransactionType *uint32
	TransactionTime *uint64
}

func (t *UpdateTransactionRequest) GetUserID() string {
	if t != nil && t.UserID != nil {
		return *t.UserID
	}
	return ""
}

func (t *UpdateTransactionRequest) GetTransactionID() string {
	if t != nil && t.TransactionID != nil {
		return *t.TransactionID
	}
	return ""
}

func (t *UpdateTransactionRequest) GetCategoryID() string {
	if t != nil && t.CategoryID != nil {
		return *t.CategoryID
	}
	return ""
}

func (t *UpdateTransactionRequest) GetAmount() string {
	if t != nil && t.Amount != nil {
		return *t.Amount
	}
	return ""
}

func (t *UpdateTransactionRequest) GetNote() string {
	if t != nil && t.Note != nil {
		return *t.Note
	}
	return ""
}

func (t *UpdateTransactionRequest) GetTransactionType() uint32 {
	if t != nil && t.TransactionType != nil {
		return *t.TransactionType
	}
	return 0
}

func (t *UpdateTransactionRequest) GetTransactionTime() uint64 {
	if t != nil && t.TransactionTime != nil {
		return *t.TransactionTime
	}
	return 0
}

func (m *UpdateTransactionRequest) ToTransactionEntity() *entity.Transaction {
	t := &entity.Transaction{
		CategoryID:      m.CategoryID,
		TransactionType: m.TransactionType,
		TransactionTime: m.TransactionTime,
		Note:            m.Note,
	}

	t.SetAmount(m.GetAmount())

	return t
}

func (m *UpdateTransactionRequest) ToGetTransactionRequest() *GetTransactionRequest {
	return &GetTransactionRequest{
		TransactionID: m.TransactionID,
	}
}

func (m *UpdateTransactionRequest) ToGetCategoryRequest() *category.GetCategoryRequest {
	return &category.GetCategoryRequest{
		CategoryID: m.CategoryID,
	}
}

func (m *UpdateTransactionRequest) ToTransactionFilter() *repo.TransactionFilter {
	return &repo.TransactionFilter{
		UserID:        m.UserID,
		TransactionID: m.TransactionID,
	}
}

type UpdateTransactionResponse struct {
	Transaction *Transaction
}

func (m *UpdateTransactionResponse) GetTransaction() *Transaction {
	if m != nil && m.Transaction != nil {
		return m.Transaction
	}
	return nil
}

type AggrTransactionsRequest struct {
	UserID           *string
	TransactionTime  *common.UInt64Filter
	CategoryIDs      []string
	TransactionTypes []uint32
}

func (t *AggrTransactionsRequest) GetUserID() string {
	if t != nil && t.UserID != nil {
		return *t.UserID
	}
	return ""
}

func (m *AggrTransactionsRequest) GetTransactionTime() *common.UInt64Filter {
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
		tt = new(common.UInt64Filter)
	}

	return &repo.TransactionFilter{
		UserID:             goutil.String(userID),
		CategoryIDs:        m.CategoryIDs,
		TransactionTypes:   m.TransactionTypes,
		TransactionTimeGte: tt.Gte,
		TransactionTimeLte: tt.Lte,
	}
}

func (m *AggrTransactionsRequest) ToGetCategoriesRequest() *category.GetCategoriesRequest {
	return &category.GetCategoriesRequest{
		CategoryIDs: m.CategoryIDs,
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
