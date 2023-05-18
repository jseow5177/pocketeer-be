package mongo

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo/model"
	"github.com/jseow5177/pockteer-be/entity"
)

const transactionCollName = "transaction"

type transactionMongo struct {
	mColl *MongoColl
}

func NewTransactionMongo(mongo *Mongo) repo.TransactionRepo {
	return &transactionMongo{
		mColl: NewMongoColl(mongo, transactionCollName),
	}
}

func (m *transactionMongo) Create(ctx context.Context, t *entity.Transaction) (string, error) {
	tm := model.ToTransactionModel(t)

	id, err := m.mColl.create(ctx, tm)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (m *transactionMongo) Update(ctx context.Context, tf *repo.TransactionFilter, t *entity.Transaction) error {
	tm := model.ToTransactionModel(t)
	if err := m.mColl.update(ctx, tf, tm); err != nil {
		return err
	}

	return nil
}

func (m *transactionMongo) Get(ctx context.Context, tf *repo.TransactionFilter) (*entity.Transaction, error) {
	t := new(model.Transaction)
	if err := m.mColl.get(ctx, tf, &t); err != nil {
		return nil, err
	}

	return model.ToTransactionEntity(t), nil
}

func (m *transactionMongo) GetMany(ctx context.Context, tf *repo.TransactionFilter) ([]*entity.Transaction, error) {
	res, err := m.mColl.getMany(ctx, tf, tf.Paging, new(model.Transaction))
	if err != nil {
		return nil, err
	}

	ets := make([]*entity.Transaction, 0, len(res))
	for _, r := range res {
		ets = append(ets, model.ToTransactionEntity(r.(*model.Transaction)))
	}

	return ets, nil
}
