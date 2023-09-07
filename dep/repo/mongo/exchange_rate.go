package mongo

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo/model"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/pkg/mongoutil"
)

const exchangeRateCollName = "exchange_rate"

type exchangeRateMongo struct {
	mColl *MongoColl
}

func NewExchangeRateMongo(mongo *Mongo) repo.ExchangeRateRepo {
	return &exchangeRateMongo{
		mColl: NewMongoColl(mongo, exchangeRateCollName),
	}
}

func (m *exchangeRateMongo) GetMany(ctx context.Context, erf *repo.ExchangeRateFilter) ([]*entity.ExchangeRate, error) {
	f := mongoutil.BuildFilter(erf)

	res, err := m.mColl.getMany(ctx, new(model.ExchangeRate), nil, f)
	if err != nil {
		return nil, err
	}

	ers := make([]*entity.ExchangeRate, 0, len(res))
	for _, r := range res {
		ers = append(ers, model.ToExchangeRateEntity(r.(*model.ExchangeRate)))
	}

	return ers, nil
}

func (m *exchangeRateMongo) CreateMany(ctx context.Context, ers []*entity.ExchangeRate) ([]string, error) {
	erms := make([]interface{}, 0)
	for _, er := range ers {
		erms = append(erms, model.ToExchangeRateModelFromEntity(er))
	}

	ids, err := m.mColl.createMany(ctx, erms)
	if err != nil {
		return nil, err
	}

	for i, er := range ers {
		er.SetExchangeRateID(goutil.String(ids[i]))
	}

	return ids, nil
}
