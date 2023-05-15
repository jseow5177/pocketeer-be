package mongo

import (
	"context"

	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/dep/repo"
)

const budgetCollName = "budget"

type budgetMongo struct {
	mColl *MongoColl
}

func NewBudgetMongo(mongo *Mongo) repo.BudgetConfigRepo {
	return &budgetMongo{
		mColl: NewMongoColl(mongo, budgetCollName),
	}
}

func (m *budgetMongo) Get(
	ctx context.Context,
	filter *repo.BudgetConfigFilter,
) (*entity.BudgetConfig, error) {
	return nil, nil
}

func (m *budgetMongo) GetMany(
	ctx context.Context,
	filter *repo.BudgetConfigFilter,
) ([]*entity.BudgetConfig, error) {
	return nil, nil
}

func (m *budgetMongo) Update(
	ctx context.Context,
	budgetConfig *entity.BudgetConfig,
) (*entity.BudgetConfig, error) {
	return nil, nil
}
