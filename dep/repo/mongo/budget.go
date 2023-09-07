package mongo

import (
	"context"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo/model"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/filter"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/pkg/mongoutil"
	"github.com/jseow5177/pockteer-be/util"
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

func (m *budgetMongo) CreateMany(ctx context.Context, bs []*entity.Budget) ([]string, error) {
	bms := make([]interface{}, 0)
	for _, b := range bs {
		bms = append(bms, model.ToBudgetModelFromEntity(b))
	}

	ids, err := m.mColl.createMany(ctx, bms)
	if err != nil {
		return nil, err
	}

	for i, b := range bs {
		b.SetBudgetID(goutil.String(ids[i]))
	}

	return ids, nil
}

func (m *budgetMongo) Get(ctx context.Context, f *repo.GetBudgetFilter) (*entity.Budget, error) {
	bq, err := m.toGetBudgetQuery(f)
	if err != nil {
		return nil, err
	}

	bs, err := m.GetMany(ctx, bq)
	if err != nil {
		return nil, err
	}

	var b *entity.Budget
	if len(bs) > 0 && !bs[0].IsDeleted() {
		b = bs[0]
	}

	if b == nil {
		return nil, repo.ErrBudgetNotFound
	}

	return b, nil
}

func (m *budgetMongo) GetMany(ctx context.Context, bq *repo.BudgetQuery) ([]*entity.Budget, error) {
	q, err := mongoutil.BuildQuery(bq)
	if err != nil {
		return nil, err
	}

	res, err := m.mColl.getMany(ctx, new(model.Budget), bq.Paging, q)
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

func (m *budgetMongo) DeleteMany(ctx context.Context, f *repo.BudgetFilter) error {
	return m.mColl.deleteMany(ctx, f)
}

func (m *budgetMongo) Delete(ctx context.Context, f *repo.DeleteBudgetFilter) error {
	startDate, endDate, err := entity.GetBudgetStartEnd(
		f.GetBudgetDate(),
		f.GetBudgetType(),
		f.GetBudgetType(),
	)
	if err != nil {
		return err
	}

	dummyBudget, err := entity.NewBudget(
		f.GetUserID(),
		f.GetCategoryID(),
		entity.WithBudgetAmount(goutil.Float64(0)),
		entity.WithBudgetType(f.BudgetType),
		entity.WithBudgetStatus(goutil.Uint32(uint32(entity.BudgetStatusDeleted))),
		entity.WithBudgetStartDate(goutil.Uint64(startDate)),
		entity.WithBudgetEndDate(goutil.Uint64(endDate)),
	)
	if err != nil {
		return err
	}

	if _, err := m.Create(ctx, dummyBudget); err != nil {
		return err
	}

	return nil
}

func (m *budgetMongo) toGetBudgetQuery(f *repo.GetBudgetFilter) (*repo.BudgetQuery, error) {
	t, err := util.ParseDateToInt(f.GetBudgetDate())
	if err != nil {
		return nil, err
	}

	return &repo.BudgetQuery{
		Queries: []*repo.BudgetQuery{
			{
				Filters: []*repo.BudgetFilter{
					{
						StartDateLte: goutil.Uint64(t),
						EndDateGte:   goutil.Uint64(t),
					},
					{
						StartDateLte: goutil.Uint64(t),
						EndDate:      goutil.Uint64(0),
					},
					{
						StartDate: goutil.Uint64(0),
						EndDate:   goutil.Uint64(0),
					},
				},
				Op: filter.Or,
			},
			{
				Filters: []*repo.BudgetFilter{
					{
						UserID:     goutil.String(f.GetUserID()),
						CategoryID: goutil.String(f.GetCategoryID()),
					},
				},
			},
		},
		Op: filter.And,
		Paging: &repo.Paging{
			Limit: goutil.Uint32(1),
			Sorts: []filter.Sort{
				&repo.Sort{
					Field: goutil.String("update_time"),
					Order: goutil.String(config.OrderDesc),
				},
			},
		},
	}, nil
}
