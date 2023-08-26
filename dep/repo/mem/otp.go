package mem

import (
	"context"
	"fmt"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mem/model"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/patrickmn/go-cache"
)

const (
	otpKeyPrefix = "otp:"
)

type otpMemCache struct {
	memCache *MemCache
}

func NewOTPMemCache(cfg *config.MemCache) (repo.OTPRepo, error) {
	memCache, err := NewMemCache(cfg)
	if err != nil {
		return nil, err
	}
	return &otpMemCache{
		memCache: memCache,
	}, nil
}

func (mc *otpMemCache) Get(ctx context.Context, of *repo.OTPFilter) (*entity.OTP, error) {
	key := mc.getKey(of.GetEmail())

	v, err := mc.memCache.Get(key)
	if err != nil {
		if err == ErrNotFound {
			return nil, repo.ErrOTPNotFound
		}
		return nil, err
	}

	om, ok := v.(*model.OTP)
	if !ok {
		return nil, repo.ErrInvalidOTP
	}

	return model.ToOTPEntity(om)
}

func (mc *otpMemCache) Set(ctx context.Context, email string, otp *entity.OTP) {
	om := model.ToOTPModelFromEntity(otp)
	mc.memCache.Set(mc.getKey(email), om, cache.DefaultExpiration)
}

func (mc *otpMemCache) getKey(email string) string {
	return fmt.Sprintf("%s:%s", otpKeyPrefix, email)
}
