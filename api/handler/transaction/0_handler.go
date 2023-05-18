package transaction

import (
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/usecase/category"
)

type transactionHandler struct {
	categoryUseCase category.UseCase
	transactionRepo repo.TransactionRepo
}

func NewTransactionHandler(categoryUseCase category.UseCase, transactionRepo repo.TransactionRepo) *transactionHandler {
	return &transactionHandler{
		categoryUseCase: categoryUseCase,
		transactionRepo: transactionRepo,
	}
}
