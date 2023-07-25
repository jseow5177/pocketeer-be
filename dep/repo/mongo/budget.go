package mongo

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo/model"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

const budgetCollName = "budget"

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

func (m *budgetMongo) Get(ctx context.Context, bf *repo.BudgetFilter) (*entity.Budget, error) {
	return nil, nil
}

func (m *budgetMongo) GetMany(ctx context.Context, bf *repo.BudgetFilter) ([]*entity.Budget, error) {
	return nil, nil
}
