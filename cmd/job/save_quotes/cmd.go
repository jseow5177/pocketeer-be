package job

import (
	"context"
	"time"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/dep/api"
	"github.com/jseow5177/pockteer-be/dep/api/finnhub"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/pkg/logger"
	"github.com/rs/zerolog/log"
)

type JobConfig struct{}

type SaveQuotes struct {
	cfg   *config.Config
	mongo *mongo.Mongo

	securityAPI  api.SecurityAPI
	holdingRepo  repo.HoldingRepo
	securityRepo repo.SecurityRepo
}

func (c *SaveQuotes) initFlags() {}

func (c *SaveQuotes) Init(ctx context.Context) (context.Context, error) {
	var err error

	c.initFlags()

	c.cfg = config.NewConfig()

	// init logger
	ctx = logger.InitZeroLog(ctx, c.cfg.Server.LogLevel)

	// init mongo
	c.mongo, err = mongo.NewMongo(ctx, c.cfg.Mongo)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to init mongo client, err: %v", err)
		return ctx, err
	}
	defer func() {
		if err != nil {
			_ = c.mongo.Close(ctx)
		}
	}()

	// init repos
	c.holdingRepo = mongo.NewHoldingMongo(c.mongo)
	c.securityRepo = mongo.NewSecurityMongo(c.mongo)

	// init apis
	c.securityAPI = finnhub.NewFinnHubMgr(c.cfg.FinnHub)

	return ctx, nil
}

func (c *SaveQuotes) Run(ctx context.Context) error {
	// get all user default holdings
	hs, err := c.holdingRepo.GetMany(ctx, &repo.HoldingFilter{
		HoldingType: goutil.Uint32(uint32(entity.HoldingTypeDefault)),
	})
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get default holdings, err: %v", err)
		return err
	}

	var (
		backOffMs  = 300
		retryCount = 10
		count      = 0
		dupHs      = make(map[string]bool)
	)
	for _, h := range hs {
		if dupHs[h.GetSymbol()] {
			continue
		}

		// rate limit, 60 API calls per minute
		if count >= 60 {
			time.Sleep(time.Minute)
		}

		var q *entity.Quote
		if err = goutil.SyncRetry(ctx, func(ctx context.Context) error {
			q, err = c.securityAPI.GetLatestQuote(ctx, &api.SecurityFilter{
				Symbol: h.Symbol,
			})
			if err != nil {
				return err
			}
			return nil
		}, retryCount, backOffMs); err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get security quote, symbol: %v, err: %v", h.GetSymbol(), err)
			return err
		}

		count++
		dupHs[h.GetSymbol()] = true

		// update security quote
		su := entity.NewSecurityUpdate(
			entity.WithUpdateSecurityQuote(q),
		)
		if err = c.securityRepo.Update(ctx, &repo.SecurityFilter{
			Symbol: h.Symbol,
		}, su); err != nil {
			log.Ctx(ctx).Error().Msgf("fail to update security quote, symbol: %v, err: %v", h.GetSymbol(), err)
			return err
		}
	}

	return nil
}

func (c *SaveQuotes) Clean(ctx context.Context) error {
	return nil
}
