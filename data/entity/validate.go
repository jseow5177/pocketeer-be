package entity

import "errors"

var (
	ErrTransactionType = errors.New("invalid transaction type")
)

func CheckTransactionType(transactionType uint32) error {
	if _, ok := TransactionTypes[transactionType]; ok {
		return nil
	}
	return ErrTransactionType
}
