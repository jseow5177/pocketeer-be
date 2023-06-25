package transaction

import (
	"github.com/jseow5177/pockteer-be/usecase/transaction"
)

type transactionHandler struct {
	transactionUseCase transaction.UseCase
}

func NewTransactionHandler(transactionUseCase transaction.UseCase) *transactionHandler {
	return &transactionHandler{
		transactionUseCase,
	}
}
