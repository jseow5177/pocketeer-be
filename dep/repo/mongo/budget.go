package mongo

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo/model"
	"github.com/jseow5177/pockteer-be/entity"
)

const (
	idBsonField = "_id"

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

func (m *budgetMongo) GetMany(
	ctx context.Context,
	bf *repo.BudgetFilter,
) ([]*entity.Budget, error) {
	res, err := m.mColl.getMany(ctx, bf, nil, new(model.Budget))
	if err != nil {
		return nil, err
	}

	entities := make([]*entity.Budget, len(res))
	for idx, r := range res {
		entities[idx] = model.ToBudgetEntity(r.(*model.Budget))
	}

	return entities, nil
}

func (m *budgetMongo) Set(ctx context.Context, budgets []*entity.Budget) error {
	models := model.ToBudgetModels(budgets)

	interfaces := make([]interface{}, len(models))
	for idx, model := range models {
		interfaces[idx] = *model
	}

	_, err := m.mColl.upsertMany(ctx, idBsonField, interfaces)
	if err != nil {
		return err
	}

	return nil
}
