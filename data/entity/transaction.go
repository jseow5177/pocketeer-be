package entity

import (
	"strconv"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
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
	CatID           *string
	Amount          *string
	TransactionType *uint32
	TransactionTime *uint64
	CreateTime      *uint64
	UpdateTime      *uint64
}

func (t *Transaction) StandardizeAmount(decimalPlaces int) error {
	af, err := strconv.ParseFloat(t.GetAmount(), 64)
	if err != nil {
		return err
	}
	t.Amount = goutil.String(goutil.FormatFloat(af, decimalPlaces))
	return nil
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
