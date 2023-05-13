package transaction

import (
	"context"

	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/data/presenter"
)

type UseCase interface {
	GetTransaction(ctx context.Context, userID string, req *presenter.GetTransactionRequest) (*entity.Transaction, error)
	GetTransactions(ctx context.Context, userID string, req *presenter.GetTransactionsRequest) ([]*entity.Transaction, error)

	CreateTransaction(ctx context.Context, userID string, req *presenter.CreateTransactionRequest) (*entity.Transaction, error)
}
