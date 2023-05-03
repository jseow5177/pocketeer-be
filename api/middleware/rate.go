package middleware

import (
	"context"
	"net/http"

	"golang.org/x/time/rate"

	"github.com/jseow5177/pockteer-be/config"
)

type RateLimiter struct {
	limiter map[string]*rate.Limiter
}

func NewRateLimiter(rls map[string]*config.RateLimit) *RateLimiter {
	if rls == nil {
		rls = make(map[string]*config.RateLimit)
	}

	limiter := make(map[string]*rate.Limiter)

	for r, rl := range rls {
		limiter[r] = rate.NewLimiter(rate.Limit(rl.RPS), int(rl.Burst))
	}

	return &RateLimiter{
		limiter: limiter,
	}
}

func RateLimit(ctx context.Context, rl *RateLimiter, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := rl.limiter[r.URL.Path]
		if limiter != nil {
			err := limiter.Wait(ctx)
			if err != nil {
				w.WriteHeader(http.StatusTooManyRequests)
			}
		}
		handler.ServeHTTP(w, r)
	})
}
