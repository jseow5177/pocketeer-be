package transaction

import (
	"context"

	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/data/presenter"
	"github.com/jseow5177/pockteer-be/dep/repo"
)

type TransactionUseCase struct {
	categoryRepo    repo.CategoryRepo
	transactionRepo repo.TransactionRepo
}

func NewTransactionUseCase(categoryRepo repo.CategoryRepo, transactionRepo repo.TransactionRepo) UseCase {
	return &TransactionUseCase{
		categoryRepo:    categoryRepo,
		transactionRepo: transactionRepo,
	}
}

func (uc *TransactionUseCase) GetTransaction(ctx context.Context, userID string, req *presenter.GetTransactionRequest) (*entity.Transaction, error) {
	return nil, nil
}

func (uc *TransactionUseCase) CreateTransaction(ctx context.Context, userID string, req *presenter.CreateTransactionRequest) (*entity.Transaction, error) {
	return nil, nil
}
