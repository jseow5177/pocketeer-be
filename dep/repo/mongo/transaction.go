package mongo

import (
	"context"
	"fmt"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo/model"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/pkg/mongoutil"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	transactionCollName = "transaction"

	aggrSumAmount = "sumAmount"
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
	tm := model.ToTransactionModelFromUpdate(tu)
	if err := m.mColl.update(ctx, tf, tm); err != nil {
		return err
	}

	return nil
}

func (m *transactionMongo) Get(ctx context.Context, tf *repo.TransactionFilter) (*entity.Transaction, error) {
	t := new(model.Transaction)
	if err := m.mColl.get(ctx, tf, &t); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repo.ErrTransactionNotFound
		}
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

func (m *transactionMongo) CalcTotalAmount(ctx context.Context, groupBy string, tf *repo.TransactionFilter) ([]*repo.TransactionAggr, error) {
	sumAmountAggr := mongoutil.NewAggr(aggrSumAmount, mongoutil.AggrSum, &mongoutil.AggrOpt{
		Field: "amount",
	})
	aggrRes, err := m.mColl.aggr(ctx, tf, groupBy, sumAmountAggr)
	if err != nil {
		return nil, err
	}

	aggrs := make([]*repo.TransactionAggr, 0)
	for _, ag := range aggrRes {
		sumAmount := ag[aggrSumAmount]

		aggrs = append(aggrs, &repo.TransactionAggr{
			GroupBy:     goutil.String(fmt.Sprint(ag[aggrGroupByField])),
			TotalAmount: goutil.Float64(mongoutil.ToFloat64(sumAmount)),
		})
	}

	return aggrs, nil
}
