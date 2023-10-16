package mongo

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jseow5177/pockteer-be/dep/api"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo/model"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/pkg/mongoutil"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const quoteCollName = "quote"

type quoteMongo struct {
	mu          sync.RWMutex
	dLock       *goutil.DLock
	mColl       *MongoColl
	securityAPI api.SecurityAPI

	quotes map[string]*entity.Quote
}

func NewQuoteMongo(ctx context.Context, mongo *Mongo, securityAPI api.SecurityAPI) (repo.QuoteRepo, error) {
	qm := &quoteMongo{
		dLock:       goutil.NewDLock(),
		mColl:       NewMongoColl(mongo, quoteCollName),
		securityAPI: securityAPI,
	}

	qs, err := qm.Load(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to load quotes mem, err: %v", err)
	}

	qm.mu.Lock()
	qm.quotes = qs
	qm.mu.Unlock()

	log.Ctx(ctx).Info().Msgf("init %v quotes in mem", len(qs))

	go func() {
		interval := 15 * time.Minute

		timer := time.NewTimer(interval)
		defer timer.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Ctx(ctx).Info().Msg("context done")
				return
			case <-timer.C:
				qs, err := qm.Load(ctx)
				if err != nil {
					log.Ctx(ctx).Error().Msgf("fail to load quotes mem, err: %v", err)
				} else {
					qm.mu.Lock()
					qm.quotes = qs
					qm.mu.Unlock()
					log.Ctx(ctx).Info().Msgf("reload quotes mem, count: %v", len(qs))
				}
			}
			timer.Reset(interval)
		}
	}()

	return qm, nil
}

func (m *quoteMongo) Load(ctx context.Context) (map[string]*entity.Quote, error) {
	var (
		page  = uint32(1)
		limit = uint32(1000)
	)

	p := &repo.Paging{
		Limit: goutil.Uint32(limit),
		Page:  goutil.Uint32(page),
	}

	quotes := make(map[string]*entity.Quote, 0)
	for {
		qf := repo.NewQuoteFilter(repo.WithQuoteRatePaging(p))
		qs, err := m.GetMany(ctx, qf)
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to scan quotes, err: %v", err)
			return nil, err
		}

		for _, q := range qs {
			quotes[q.GetSymbol()] = q
		}

		if uint32(len(qs)) == limit {
			page += uint32(1)
			p.Page = goutil.Uint32(page)
		} else {
			break
		}
	}

	return quotes, nil
}

func (m *quoteMongo) Get(ctx context.Context, qf *repo.QuoteFilter) (*entity.Quote, error) {
	m.mu.RLock()

	symbol := qf.GetSymbol()

	// Get from cache
	q, ok := m.quotes[symbol]
	if ok {
		m.mu.RUnlock()
		return q, nil
	}

	m.mu.RUnlock()

	// Try to get lock to access third-party API.
	// If cannot get lock, call API again. Quote may already be in cache.
	if locked := m.dLock.TryLock(symbol); !locked {
		time.Sleep(100 * time.Millisecond)
		return m.Get(ctx, qf)
	}

	defer func() {
		m.dLock.Unlock(symbol)
	}()

	// Get from third-party API
	q, err := m.securityAPI.GetLatestQuote(ctx, &api.SecurityFilter{
		Symbol: qf.Symbol,
	})
	if err != nil {
		return nil, err
	}

	// Set to cache
	m.mu.Lock()
	m.quotes[symbol] = q
	m.mu.Unlock()

	return q, nil
}

func (m *quoteMongo) Upsert(ctx context.Context, qf *repo.QuoteFilter, q *entity.Quote) error {
	qm := model.ToQuoteModelFromEntity(q)

	f := mongoutil.BuildFilter(qf)

	return m.mColl.update(ctx, f, qm, options.Update().SetUpsert(true))
}

func (m *quoteMongo) GetMany(ctx context.Context, qf *repo.QuoteFilter) ([]*entity.Quote, error) {
	f := mongoutil.BuildFilter(qf)

	res, err := m.mColl.getMany(ctx, new(model.Quote), qf.Paging, f)
	if err != nil {
		return nil, err
	}

	qs := make([]*entity.Quote, 0, len(res))
	for _, r := range res {
		qs = append(qs, model.ToQuoteEntity(r.(*model.Quote)))
	}

	return qs, nil
}
