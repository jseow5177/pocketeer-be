package presenter

import (
	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/common"
	"github.com/jseow5177/pockteer-be/usecase/transaction"
)

type Transaction struct {
	TransactionID   *string   `json:"transaction_id,omitempty"`
	Category        *Category `json:"category,omitempty"`
	Amount          *string   `json:"amount,omitempty"`
	Note            *string   `json:"note,omitempty"`
	TransactionType *uint32   `json:"transaction_type,omitempty"`
	TransactionTime *uint64   `json:"transaction_time,omitempty"`
	CreateTime      *uint64   `json:"create_time,omitempty"`
	UpdateTime      *uint64   `json:"update_time,omitempty"`
}

func (t *Transaction) GetTransactionID() string {
	if t != nil && t.TransactionID != nil {
		return *t.TransactionID
	}
	return ""
}

func (t *Transaction) GetCategory() *Category {
	if t != nil && t.Category != nil {
		return t.Category
	}
	return nil
}

func (t *Transaction) GetAmount() string {
	if t != nil && t.Amount != nil {
		return *t.Amount
	}
	return ""
}

func (t *Transaction) GetNote() string {
	if t != nil && t.Note != nil {
		return *t.Note
	}
	return ""
}

func (t *Transaction) GetTransactionType() uint32 {
	if t != nil && t.TransactionType != nil {
		return *t.TransactionType
	}
	return 0
}

func (t *Transaction) GetTransactionTime() uint64 {
	if t != nil && t.TransactionTime != nil {
		return *t.TransactionTime
	}
	return 0
}

func (t *Transaction) GetCreateTime() uint64 {
	if t != nil && t.CreateTime != nil {
		return *t.CreateTime
	}
	return 0
}

func (t *Transaction) GetUpdateTime() uint64 {
	if t != nil && t.UpdateTime != nil {
		return *t.UpdateTime
	}
	return 0
}

type CreateTransactionRequest struct {
	CategoryID      *string `json:"category_id,omitempty"`
	Amount          *string `json:"amount,omitempty"`
	TransactionType *uint32 `json:"transaction_type,omitempty"`
	TransactionTime *uint64 `json:"transaction_time,omitempty"`
	Note            *string `json:"note,omitempty"`
}

func (m *CreateTransactionRequest) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *CreateTransactionRequest) GetAmount() string {
	if m != nil && m.Amount != nil {
		return *m.Amount
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

func (m *CreateTransactionRequest) GetNote() string {
	if m != nil && m.Note != nil {
		return *m.Note
	}
	return ""
}

func (m *CreateTransactionRequest) ToUseCaseReq(userID string) *transaction.CreateTransactionRequest {
	return &transaction.CreateTransactionRequest{
		UserID:          goutil.String(userID),
		CategoryID:      m.CategoryID,
		Amount:          m.Amount,
		TransactionType: m.TransactionType,
		TransactionTime: m.TransactionTime,
		Note:            m.Note,
	}
}

type CreateTransactionResponse struct {
	Transaction *Transaction `json:"transaction,omitempty"`
}

func (m *CreateTransactionResponse) GetTransaction() *Transaction {
	if m != nil && m.Transaction != nil {
		return m.Transaction
	}
	return nil
}

func (m *CreateTransactionResponse) Set(useCaseRes *transaction.CreateTransactionResponse) {
	m.Transaction = toTransaction(useCaseRes.Transaction, useCaseRes.Category)
}

type GetTransactionRequest struct {
	TransactionID *string `json:"transaction_id,omitempty"`
}

func (m *GetTransactionRequest) GetTransactionID() string {
	if m != nil && m.TransactionID != nil {
		return *m.TransactionID
	}
	return ""
}

func (m *GetTransactionRequest) ToUseCaseReq(userID string) *transaction.GetTransactionRequest {
	return &transaction.GetTransactionRequest{
		UserID:        goutil.String(userID),
		TransactionID: m.TransactionID,
	}
}

type GetTransactionResponse struct {
	Transaction *Transaction `json:"transaction,omitempty"`
}

func (m *GetTransactionResponse) GetTransaction() *Transaction {
	if m != nil && m.Transaction != nil {
		return m.Transaction
	}
	return nil
}

func (m *GetTransactionResponse) Set(useCaseRes *transaction.GetTransactionResponse) {
	m.Transaction = toTransaction(useCaseRes.Transaction, useCaseRes.Category)
}

type GetTransactionsRequest struct {
	CategoryID      *string       `json:"category_id,omitempty"`
	TransactionType *uint32       `json:"transaction_type,omitempty"`
	TransactionTime *UInt64Filter `json:"transaction_time,omitempty"`
	Paging          *Paging       `json:"paging,omitempty"`
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

func (m *GetTransactionsRequest) GetTransactionTime() *UInt64Filter {
	if m != nil && m.TransactionTime != nil {
		return m.TransactionTime
	}
	return nil
}

func (m *GetTransactionsRequest) GetPaging() *Paging {
	if m != nil && m.Paging != nil {
		return m.Paging
	}
	return nil
}

func (m *GetTransactionsRequest) ToUseCaseReq(userID string) *transaction.GetTransactionsRequest {
	paging := m.Paging
	if paging == nil {
		paging = new(Paging)
	}

	if paging.Limit == nil {
		paging.Limit = goutil.Uint32(config.DefaultPagingLimit)
	}

	if paging.Page == nil {
		paging.Page = goutil.Uint32(config.MinPagingPage)
	}

	tt := m.TransactionTime
	if tt == nil {
		tt = new(UInt64Filter)
	}

	return &transaction.GetTransactionsRequest{
		UserID:          goutil.String(userID),
		CategoryID:      m.CategoryID,
		TransactionType: m.TransactionType,
		Paging: &common.Paging{
			Limit: paging.Limit,
			Page:  paging.Page,
		},
		TransactionTime: &common.UInt64Filter{
			Gte: tt.Gte,
			Lte: tt.Lte,
		},
	}
}

type GetTransactionsResponse struct {
	Transactions []*Transaction `json:"transactions,omitempty"`
	Paging       *Paging        `json:"paging,omitempty"`
}

func (m *GetTransactionsResponse) GetTransactions() []*Transaction {
	if m != nil && m.Transactions != nil {
		return m.Transactions
	}
	return nil
}

func (m *GetTransactionsResponse) GetPaging() *Paging {
	if m != nil && m.Paging != nil {
		return m.Paging
	}
	return nil
}

func (m *GetTransactionsResponse) Set(useCaseRes *transaction.GetTransactionsResponse) {
	ts := make([]*Transaction, 0)
	for _, twc := range useCaseRes.TransactionsWithCategory {
		ts = append(ts, toTransaction(twc.Transaction, twc.Category))
	}
	m.Transactions = ts
	m.Paging = toPaging(useCaseRes.Paging)
}

type UpdateTransactionRequest struct {
	TransactionID   *string `json:"transaction_id,omitempty"`
	CategoryID      *string `json:"category_id,omitempty"`
	Amount          *string `json:"amount,omitempty"`
	Note            *string `json:"note,omitempty"`
	TransactionType *uint32 `json:"transaction_type,omitempty"`
	TransactionTime *uint64 `json:"transaction_time,omitempty"`
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

func (m *UpdateTransactionRequest) ToUseCaseReq(userID string) *transaction.UpdateTransactionRequest {
	return &transaction.UpdateTransactionRequest{
		UserID:          goutil.String(userID),
		TransactionID:   m.TransactionID,
		CategoryID:      m.CategoryID,
		Note:            m.Note,
		Amount:          m.Amount,
		TransactionType: m.TransactionType,
		TransactionTime: m.TransactionTime,
	}
}

type UpdateTransactionResponse struct {
	Transaction *Transaction `json:"transaction,omitempty"`
}

func (m *UpdateTransactionResponse) GetTransaction() *Transaction {
	if m != nil && m.Transaction != nil {
		return m.Transaction
	}
	return nil
}

func (m *UpdateTransactionResponse) Set(useCaseRes *transaction.UpdateTransactionResponse) {
	m.Transaction = toTransaction(useCaseRes.TransactionWithCategory.Transaction, useCaseRes.TransactionWithCategory.Category)
}

type AggrTransactionsRequest struct {
	TransactionTime  *UInt64Filter `json:"transaction_time,omitempty"`
	CategoryIDs      []string      `json:"category_ids,omitempty"`
	TransactionTypes []uint32      `json:"transaction_types,omitempty"`
}

func (m *AggrTransactionsRequest) GetTransactionTime() *UInt64Filter {
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

func (m *AggrTransactionsRequest) ToUseCaseReq(userID string) *transaction.AggrTransactionsRequest {
	tt := m.TransactionTime
	if tt == nil {
		tt = new(UInt64Filter)
	}

	return &transaction.AggrTransactionsRequest{
		UserID:           goutil.String(userID),
		TransactionTypes: m.TransactionTypes,
		CategoryIDs:      m.CategoryIDs,
		TransactionTime: &common.UInt64Filter{
			Gte: tt.Gte,
			Lte: tt.Lte,
		},
	}
}

type Aggr struct {
	Sum *float64 `json:"sum,omitempty"`
}

type AggrTransactionsResponse struct {
	Results map[string]*Aggr `json:"results,omitempty"`
}

func (m *AggrTransactionsResponse) GetResults() map[string]*Aggr {
	if m != nil && m.Results != nil {
		return m.Results
	}
	return nil
}

func (m *AggrTransactionsResponse) Set(useCaseRes *transaction.AggrTransactionsResponse) {
	res := make(map[string]*Aggr)
	for k, v := range useCaseRes.Results {
		res[k] = toAggr(v)
	}
	m.Results = res
}
