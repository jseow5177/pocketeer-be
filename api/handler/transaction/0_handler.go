package transaction

import (
	"github.com/jseow5177/pockteer-be/usecase/aggr"
	"github.com/jseow5177/pockteer-be/usecase/transaction"
)

type transactionHandler struct {
	transactionUseCase transaction.UseCase
	aggrUseCase        aggr.UseCase
}

func NewTransactionHandler(transactionUseCase transaction.UseCase, aggrUseCase aggr.UseCase) *transactionHandler {
	return &transactionHandler{
		transactionUseCase,
		aggrUseCase,
	}
}
