package presenter

import (
	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

type Transaction struct {
	TransactionID   *string `json:"transaction_id,omitempty"`
	UserID          *string `json:"user_id,omitempty"`
	CatID           *string `json:"cat_id,omitempty"`
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

func (t *Transaction) GetUserID() string {
	if t != nil && t.UserID != nil {
		return *t.UserID
	}
	return ""
}

func (t *Transaction) GetCatID() string {
	if t != nil && t.CatID != nil {
		return *t.CatID
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
		UserID:          t.UserID,
		CatID:           t.CatID,
		Amount:          t.Amount,
		TransactionType: t.TransactionType,
		TransactionTime: t.TransactionTime,
		CreateTime:      t.CreateTime,
		UpdateTime:      t.UpdateTime,
	}
}

type CreateTransactionRequest struct {
	CatID           *string `json:"cat_id,omitempty"`
	Amount          *string `json:"amount,omitempty"`
	TransactionType *uint32 `json:"transaction_type,omitempty"`
	TransactionTime *uint64 `json:"transaction_time,omitempty"`
}

func (m *CreateTransactionRequest) GetCatID() string {
	if m != nil && m.CatID != nil {
		return *m.CatID
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
		CatID:           m.CatID,
		Amount:          m.Amount,
		TransactionType: m.TransactionType,
		TransactionTime: m.TransactionTime,
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
