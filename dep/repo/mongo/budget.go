package mongo

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
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

func (r *budgetMongo) GetMany(ctx context.Context, req *repo.BudgetFilter) ([]*entity.Budget, error) {
	return nil, nil
}

func (r *budgetMongo) Set(ctx context.Context, budgets []*entity.Budget) error {
	return nil
}
