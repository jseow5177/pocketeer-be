package transaction

import (
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/usecase/category"
)

type TransactionHandler struct {
	categoryUseCase category.UseCase
	transactionRepo repo.TransactionRepo
}

func NewTransactionHandler(categoryUseCase category.UseCase, transactionRepo repo.TransactionRepo) *TransactionHandler {
	return &TransactionHandler{
		categoryUseCase: categoryUseCase,
		transactionRepo: transactionRepo,
	}
}
