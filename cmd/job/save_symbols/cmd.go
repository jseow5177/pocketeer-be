package job

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

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

const (
	DefaultExchange = "US"
)

type JobConfig struct {
	Exchange string
}

type SaveSymbols struct {
	jobCfg JobConfig

	cfg   *config.Config
	mongo *mongo.Mongo

	securityAPI  api.SecurityAPI
	securityRepo repo.SecurityRepo
}

func (c *SaveSymbols) initFlags() {
	flagSet := flag.NewFlagSet(fmt.Sprintf("%s %s", filepath.Base(os.Args[0]), os.Args[1]), flag.ExitOnError)

	flagSet.StringVar(&c.jobCfg.Exchange, "exchange", DefaultExchange, "exchange of symbols")
}

func (c *SaveSymbols) Init(ctx context.Context) (context.Context, error) {
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
	c.securityRepo = mongo.NewSecurityMongo(c.mongo)

	// init apis
	c.securityAPI = finnhub.NewFinnHubMgr(c.cfg.FinnHub)

	return ctx, nil
}

func (c *SaveSymbols) Run(ctx context.Context) error {
	// scan symbols from API
	ss, err := c.securityAPI.ListSymbols(ctx, &api.SecurityFilter{
		Exchange: goutil.String(c.jobCfg.Exchange),
	})
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to list symbols from api, err: %v", err)
		return err
	}

	var (
		batchSize  = 1000
		batch      = make([]*entity.Security, 0, batchSize)
		retryCount = 10
		backoffMs  = 300
		count      = 0
	)

	for _, s := range ss {
		batch = append(batch, s)
		if len(batch) < batchSize {
			continue
		}

		if err := goutil.SyncRetry(ctx, func(ctx context.Context) error {
			return c.securityRepo.CreateMany(ctx, batch)
		}, retryCount, backoffMs); err != nil {
			log.Ctx(ctx).Error().Msgf("fail to save securities to mongo, err: %v", err)
			return err
		}

		count += len(batch)

		batch = make([]*entity.Security, 0, 1000)
	}

	if len(batch) > 0 {
		if err := goutil.SyncRetry(ctx, func(ctx context.Context) error {
			return c.securityRepo.CreateMany(ctx, batch)
		}, retryCount, backoffMs); err != nil {
			log.Ctx(ctx).Error().Msgf("fail to save securities to mongo, err: %v", err)
			return err
		}

		count += len(batch)
	}

	log.Ctx(ctx).Info().Msgf("inserted %v symbols, exchange: %v", count, c.jobCfg.Exchange)

	return nil
}

func (c *SaveSymbols) Clean(ctx context.Context) error {
	return nil
}
