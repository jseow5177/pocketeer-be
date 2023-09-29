package mongo

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo/model"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/filter"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/pkg/mongoutil"
	"github.com/rs/zerolog/log"
)

const exchangeRateCollName = "exchange_rate"

type exchangeRateMongo struct {
	mu    sync.RWMutex
	mColl *MongoColl

	exchangeRates map[string][]*entity.ExchangeRate // from-to -> exchange rates
}

func NewExchangeRateMongo(ctx context.Context, mongo *Mongo) (repo.ExchangeRateRepo, error) {
	erm := &exchangeRateMongo{
		mColl: NewMongoColl(mongo, exchangeRateCollName),
	}

	ers, err := erm.Load(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to load exchange rate mem, err: %v", err)
	}

	erm.mu.Lock()
	erm.exchangeRates = ers
	erm.mu.Unlock()

	var total int
	for _, rates := range ers {
		total += len(rates)
	}
	log.Ctx(ctx).Info().Msgf("init %v exchange rates in mem", total)

	go func() {
		interval := 24 * time.Hour

		timer := time.NewTimer(interval)
		defer timer.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Ctx(ctx).Info().Msg("context done")
				return
			case <-timer.C:
				ers, err := erm.Load(ctx)
				if err != nil {
					log.Ctx(ctx).Error().Msgf("fail to load exchange rate mem, err: %v", err)
				} else {
					erm.mu.Lock()
					erm.exchangeRates = ers
					erm.mu.Unlock()
				}
			}
		}
	}()

	return erm, nil
}

func (m *exchangeRateMongo) Load(ctx context.Context) (map[string][]*entity.ExchangeRate, error) {
	var (
		page  = uint32(1)
		limit = uint32(1000)
	)

	p := &repo.Paging{
		Limit: goutil.Uint32(limit),
		Page:  goutil.Uint32(page),
		Sorts: []filter.Sort{
			&repo.Sort{
				Field: goutil.String("timestamp"),
				Order: goutil.String(config.OrderAsc),
			},
		},
	}

	exchangeRates := make([]*entity.ExchangeRate, 0)
	for {
		erf := repo.NewExchangeRateFilter(repo.WithExchangeRatePaging(p))
		ers, err := m.GetMany(ctx, erf)
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to scan exchange rates, err: %v", err)
			return nil, err
		}

		exchangeRates = append(exchangeRates, ers...)

		if uint32(len(ers)) == limit {
			page += uint32(1)
			p.Page = goutil.Uint32(page)
		} else {
			break
		}
	}

	exchangeRateGroups := make(map[string][]*entity.ExchangeRate)
	for _, er := range exchangeRates {
		k := fmt.Sprintf("%s-%s", er.GetFrom(), er.GetTo())
		exchangeRateGroups[k] = append(exchangeRateGroups[k], er)
	}

	return exchangeRateGroups, nil
}

func (m *exchangeRateMongo) GetMany(ctx context.Context, erf *repo.ExchangeRateFilter) ([]*entity.ExchangeRate, error) {
	f := mongoutil.BuildFilter(erf)

	res, err := m.mColl.getMany(ctx, new(model.ExchangeRate), erf.Paging, f)
	if err != nil {
		return nil, err
	}

	ers := make([]*entity.ExchangeRate, 0, len(res))
	for _, r := range res {
		ers = append(ers, model.ToExchangeRateEntity(r.(*model.ExchangeRate)))
	}

	return ers, nil
}

func (m *exchangeRateMongo) Get(ctx context.Context, erf *repo.ExchangeRateFilter) (*entity.ExchangeRate, error) {
	er := m.binarySearchExchangeRates(erf)
	if er == nil {
		return nil, repo.ErrExchangeRateNotFound
	}

	return er, nil
}

func (m *exchangeRateMongo) binarySearchExchangeRates(erf *repo.ExchangeRateFilter) *entity.ExchangeRate {
	m.mu.RLock()
	defer m.mu.RUnlock()

	k := fmt.Sprintf("%s-%s", erf.GetFrom(), erf.GetTo())
	ers := m.exchangeRates[k]

	index := goutil.BinarySearch(len(ers), func(index int) bool {
		er := ers[index]
		return er.GetTimestamp() <= erf.GetTimestamp()
	})

	if index != -1 {
		return ers[index]
	}

	return nil
}

func (m *exchangeRateMongo) Create(ctx context.Context, er *entity.ExchangeRate) (string, error) {
	erm := model.ToExchangeRateModelFromEntity(er)
	id, err := m.mColl.create(ctx, erm)
	if err != nil {
		return "", err
	}
	er.SetExchangeRateID(goutil.String(id))

	return id, nil
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
