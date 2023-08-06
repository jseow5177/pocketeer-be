package mem

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/dep/api"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mem/model"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/patrickmn/go-cache"
)

var (
	ErrInvalidQuote = errors.New("invalid quote in mem cache")
)

const keyPrefix = "quote"

type quoteMemCache struct {
	dLock       *goutil.DLock
	memCache    *MemCache
	securityAPI api.SecurityAPI
}

func NewQuoteMemCache(cfg *config.MemCache, securityAPI api.SecurityAPI) (repo.QuoteRepo, error) {
	memCache, err := NewMemCache(cfg)
	if err != nil {
		return nil, err
	}
	return &quoteMemCache{
		dLock:       goutil.NewDLock(),
		memCache:    memCache,
		securityAPI: securityAPI,
	}, nil
}

func (mc *quoteMemCache) Get(ctx context.Context, qf *repo.QuoteFilter) (*entity.Quote, error) {
	key := mc.getKey(qf.GetSymbol())

	v, err := mc.memCache.Get(key)
	if err != nil && err != ErrNotFound {
		return nil, err
	}

	// exist in cache
	if err == nil {
		qm, ok := v.(*model.Quote)
		if !ok {
			return nil, ErrInvalidQuote
		}
		return model.ToQuoteEntity(qm), nil
	}

	// lock to prevent overloading third-party API
	if locked := mc.dLock.TryLock(key); locked {
		defer func() {
			mc.dLock.Unlock(key)
		}()

		// get from third-party API
		q, err := mc.securityAPI.GetLatestQuote(ctx, &api.SecurityFilter{
			Symbol: qf.Symbol,
		})
		if err != nil {
			return nil, err
		}

		// set back to cache
		mc.Set(ctx, qf.GetSymbol(), q)

		return q, nil
	} else {
		// try again, which may already be in cache
		time.Sleep(100 * time.Millisecond)
		return mc.Get(ctx, qf)
	}
}

func (mc *quoteMemCache) Set(ctx context.Context, symbol string, q *entity.Quote) {
	qm := model.ToQuoteModelFromEntity(q)
	mc.memCache.Set(mc.getKey(symbol), qm, cache.DefaultExpiration)
}

func (mc *quoteMemCache) getKey(symbol string) string {
	return fmt.Sprintf("%s:%s", keyPrefix, symbol)
}
