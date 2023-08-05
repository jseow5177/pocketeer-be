package mem

import (
	"errors"
	"time"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/patrickmn/go-cache"
)

var (
	ErrNotFound = errors.New("not found in mem cache")
)

type MemCache struct {
	client *cache.Cache
}

func NewMemCache(cfg *config.MemCache) (*MemCache, error) {
	var err error

	et := cache.DefaultExpiration
	if cfg.ExpiryTime != "" {
		et, err = time.ParseDuration(cfg.ExpiryTime)
		if err != nil {
			return nil, err
		}
	}

	var cui time.Duration
	if cfg.CleanUpInterval != "" {
		cui, err = time.ParseDuration(cfg.CleanUpInterval)
		if err != nil {
			return nil, err
		}
	}

	return &MemCache{
		client: cache.New(et, cui),
	}, nil
}

func (mc *MemCache) Set(key string, val interface{}, d time.Duration) {
	mc.client.Set(key, val, d)
}

func (mc *MemCache) Get(key string) (interface{}, error) {
	if val, found := mc.client.Get(key); found {
		return val, nil
	}
	return nil, ErrNotFound
}
