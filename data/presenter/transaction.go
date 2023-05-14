package presenter

import (
	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/pkg/filter"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

type Transaction struct {
	TransactionID   *string `json:"transaction_id,omitempty"`
	CategoryID      *string `json:"category_id,omitempty"`
	Amount          *string `json:"amount,omitempty"`
	TransactionType *uint32 `json:"transaction_type,omitempty"`
	TransactionTime *uint64 `json:"transaction_time,omitempty"`
	CreateTime      *uint64 `json:"create_time,omitempty"`
	UpdateTime      *uint64 `json:"update_time,omitempty"`
}

func (t *Transaction) GetTransactionID() string {
	if t != nil && t.TransactionID != nil {
		return *t.TransactionID
	}
	return ""
}

func (t *Transaction) GetCategoryID() string {
	if t != nil && t.CategoryID != nil {
		return *t.CategoryID
	}
	return ""
}

func (t *Transaction) GetAmount() string {
	if t != nil && t.Amount != nil {
		return *t.Amount
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

func ToTransactionPresenter(t *entity.Transaction) *Transaction {
	return &Transaction{
		TransactionID:   t.TransactionID,
		CategoryID:      t.CategoryID,
		Amount:          t.Amount,
		TransactionType: t.TransactionType,
		TransactionTime: t.TransactionTime,
		CreateTime:      t.CreateTime,
		UpdateTime:      t.UpdateTime,
	}
}

type CreateTransactionRequest struct {
	CategoryID      *string `json:"category_id,omitempty"`
	Amount          *string `json:"amount,omitempty"`
	TransactionType *uint32 `json:"transaction_type,omitempty"`
	TransactionTime *uint64 `json:"transaction_time,omitempty"`
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

func (m *CreateTransactionRequest) ToTransactionEntity(userID string) *entity.Transaction {
	return &entity.Transaction{
		UserID:          goutil.String(userID),
		CategoryID:      m.CategoryID,
		Amount:          m.Amount,
		TransactionType: m.TransactionType,
		TransactionTime: m.TransactionTime,
	}
}

func (m *CreateTransactionRequest) ToGetCategoryRequest() *GetCategoryRequest {
	return &GetCategoryRequest{
		CategoryID: m.CategoryID,
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

func (m *CreateTransactionResponse) SetTransaction(t *entity.Transaction) {
	m.Transaction = ToTransactionPresenter(t)
}

type GetTransactionRequest struct {
	TransactionID *string `json:"transaction_id"`
}

func (m *GetTransactionRequest) GetTransactionID() string {
	if m != nil && m.TransactionID != nil {
		return *m.TransactionID
	}
	return ""
}

func (m *GetTransactionRequest) ToTransactionFilter(userID string) *repo.TransactionFilter {
	return &repo.TransactionFilter{
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

func (m *GetTransactionResponse) SetTransaction(t *entity.Transaction) {
	m.Transaction = ToTransactionPresenter(t)
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

func (m *GetTransactionsRequest) ToGetCategoryRequest() *GetCategoryRequest {
	return &GetCategoryRequest{
		CategoryID: m.CategoryID,
	}
}

func (m *GetTransactionsRequest) ToTransactionFilter(userID string) *repo.TransactionFilter {
	tt := new(UInt64Filter)
	if m.TransactionTime != nil {
		tt = m.TransactionTime
	}

	paging := new(Paging)
	if m.Paging != nil {
		paging = m.Paging
	}

	return &repo.TransactionFilter{
		UserID:             goutil.String(userID),
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
	Transactions []*Transaction `json:"transactions"`
	Paging       *Paging        `json:"paging"`
}

func (m *GetTransactionsResponse) GetTransactions() []*Transaction {
	if m != nil && m.Transactions != nil {
		return m.Transactions
	}
	return nil
}

func (m *GetTransactionsResponse) SetTransactions(ets []*entity.Transaction) {
	ts := make([]*Transaction, 0)
	for _, et := range ets {
		ts = append(ts, ToTransactionPresenter(et))
	}
	m.Transactions = ts
}

func (m *GetTransactionsResponse) SetPaging(paging *Paging) {
	m.Paging = paging
}
