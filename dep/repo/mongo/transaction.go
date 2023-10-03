package mongo

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo/model"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/pkg/mongoutil"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	transactionCollName = "transaction"
)

type transactionMongo struct {
	mColl *MongoColl
}

func NewTransactionMongo(mongo *Mongo) repo.TransactionRepo {
	return &transactionMongo{
		mColl: NewMongoColl(mongo, transactionCollName),
	}
}

func (m *transactionMongo) Create(ctx context.Context, t *entity.Transaction) (string, error) {
	tm := model.ToTransactionModelFromEntity(t)
	id, err := m.mColl.create(ctx, tm)
	if err != nil {
		return "", err
	}
	t.SetTransactionID(goutil.String(id))

	return id, nil
}

func (m *transactionMongo) Update(ctx context.Context, tf *repo.TransactionFilter, tu *entity.TransactionUpdate) error {
	f := mongoutil.BuildFilter(tf)

	tm := model.ToTransactionModelFromUpdate(tu)
	if err := m.mColl.update(ctx, f, tm); err != nil {
		return err
	}

	return nil
}

func (m *transactionMongo) Get(ctx context.Context, tf *repo.TransactionFilter) (*entity.Transaction, error) {
	f := mongoutil.BuildFilter(tf)

	t := new(model.Transaction)
	if err := m.mColl.get(ctx, &t, f); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repo.ErrTransactionNotFound
		}
		return nil, err
	}

	return model.ToTransactionEntity(t)
}

func (m *transactionMongo) GetMany(ctx context.Context, tq *repo.TransactionQuery) ([]*entity.Transaction, error) {
	q, err := mongoutil.BuildQuery(tq)
	if err != nil {
		return nil, err
	}

	res, err := m.mColl.getMany(ctx, new(model.Transaction), tq.Paging, q)
	if err != nil {
		return nil, err
	}

	ets := make([]*entity.Transaction, 0, len(res))
	for _, r := range res {
		t, err := model.ToTransactionEntity(r.(*model.Transaction))
		if err != nil {
			return nil, err
		}
		ets = append(ets, t)
	}

	return ets, nil
}
