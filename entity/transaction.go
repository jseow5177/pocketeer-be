package entity

import (
	"strconv"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

const (
	DefaultAmount       = 0
	AmountDecimalPlaces = 2
)

type TransactionType uint32

const (
	TransactionTypeExpense TransactionType = 1
	TransactionTypeIncome  TransactionType = 2
)

var TransactionTypes = map[uint32]string{
	uint32(TransactionTypeExpense): "expense",
	uint32(TransactionTypeIncome):  "income",
}

type Transaction struct {
	TransactionID   *string
	UserID          *string
	CategoryID      *string
	amount          *string
	Note            *string
	TransactionType *uint32
	TransactionTime *uint64
	CreateTime      *uint64
	UpdateTime      *uint64
}

func (t *Transaction) IsNilAmount() bool {
	return t.amount == nil
}

func (t *Transaction) SetAmount(amount string) {
	af, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		af = DefaultAmount
	}
	t.amount = goutil.String(goutil.FormatFloat(af, AmountDecimalPlaces))
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

func (t *Transaction) GetCategoryID() string {
	if t != nil && t.CategoryID != nil {
		return *t.CategoryID
	}
	return ""
}

func (t *Transaction) GetAmount() string {
	if t != nil && t.amount != nil {
		return *t.amount
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
