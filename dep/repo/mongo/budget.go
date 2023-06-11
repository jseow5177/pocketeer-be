package mongo

import (
	"context"
	"time"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo/model"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	idBsonField    = "_id"
	budgetCollName = "budget"
)

type budgetMongo struct {
	mColl *MongoColl
}

func NewBudgetMongo(mongo *Mongo) repo.BudgetRepo {
	return &budgetMongo{
		mColl: NewMongoColl(mongo, budgetCollName),
	}
}

func (m *budgetMongo) Get(
	ctx context.Context,
	bf *repo.BudgetFilter,
) (*entity.Budget, error) {
	budget := new(model.Budget)

	if err := m.mColl.get(ctx, bf, &budget); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repo.ErrBudgetNotFound
		}
		return nil, err
	}

	return model.ToBudgetEntity(budget), nil
}

func (m *budgetMongo) GetMany(
	ctx context.Context,
	bf *repo.BudgetFilter,
) ([]*entity.Budget, error) {
	res, err := m.mColl.getMany(ctx, bf, nil, new(model.Budget))
	if err != nil {
		return nil, err
	}

	budgets := make([]*entity.Budget, 0, len(res))
	for _, r := range res {
		budgets = append(budgets, model.ToBudgetEntity(r.(*model.Budget)))
	}

	return budgets, nil
}

func (m *budgetMongo) Set(
	ctx context.Context,
	budget *entity.Budget,
) error {
	now := uint64(time.Now().Unix())
	budget.UpdateTime = goutil.Uint64(now)
	model := model.ToBudgetModel(budget)

	_, err := m.mColl.upsert(ctx, idBsonField, model)
	if err != nil {
		return err
	}

	return nil
}
