package transaction

import "github.com/jseow5177/pockteer-be/dep/repo"

type TransactionHandler struct {
	categoryRepo    repo.CategoryRepo
	transactionRepo repo.TransactionRepo
}

func NewTransactionHandler(categoryRepo repo.CategoryRepo, transactionRepo repo.TransactionRepo) *TransactionHandler {
	return &TransactionHandler{
		categoryRepo:    categoryRepo,
		transactionRepo: transactionRepo,
	}
}
