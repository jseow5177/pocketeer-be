package mongo

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo/model"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/pkg/mongoutil"
)

const budgetCollName = "budgetV2"

type budgetMongo struct {
	mColl *MongoColl
}

func NewBudgetMongo(mongo *Mongo) repo.BudgetRepo {
	return &budgetMongo{
		mColl: NewMongoColl(mongo, budgetCollName),
	}
}

func (m *budgetMongo) Create(ctx context.Context, b *entity.Budget) (string, error) {
	bm := model.ToBudgetModelFromEntity(b)
	id, err := m.mColl.create(ctx, bm)
	if err != nil {
		return "", err
	}
	b.SetBudgetID(goutil.String(id))

	return id, nil
}

func (m *budgetMongo) GetMany(ctx context.Context, paging *repo.Paging, bq *repo.BudgetQuery) ([]*entity.Budget, error) {
	q, err := mongoutil.BuildQuery(bq)
	if err != nil {
		return nil, err
	}

	res, err := m.mColl.getMany(ctx, new(model.Budget), paging, q)
	if err != nil {
		return nil, err
	}

	bs := make([]*entity.Budget, 0, len(res))
	for _, r := range res {
		b, err := model.ToBudgetEntity(r.(*model.Budget))
		if err != nil {
			return nil, err
		}
		bs = append(bs, b)
	}

	return bs, nil
}
