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
	"github.com/rs/zerolog/log"
)

const DefaultExchange = "US"

type JobConfig struct {
	Exchange string
}

type InitSymbols struct {
	cfg JobConfig

	mongo *mongo.Mongo

	securityAPI  api.SecurityAPI
	securityRepo repo.SecurityRepo
}

func (c *InitSymbols) initFlags() error {
	flagSet := flag.NewFlagSet(fmt.Sprintf("%s %s", filepath.Base(os.Args[0]), os.Args[1]), flag.ExitOnError)

	flagSet.StringVar(&c.cfg.Exchange, "exchange", DefaultExchange, "exchange of symbols")

	if err := flagSet.Parse(os.Args[2:]); err != nil {
		return err
	}

	return nil
}

func (c *InitSymbols) Init(ctx context.Context, cfg *config.Config) error {
	var err error

	if err = c.initFlags(); err != nil {
		return err
	}

	// init mongo
	c.mongo, err = mongo.NewMongo(ctx, cfg.Mongo)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to init mongo client, err: %v", err)
		return err
	}
	defer func() {
		if err != nil {
			_ = c.mongo.Close(ctx)
		}
	}()

	// init repos
	c.securityRepo = mongo.NewSecurityMongo(c.mongo)

	// init apis
	c.securityAPI = finnhub.NewFinnHubMgr(cfg.FinnHub)

	return nil
}

func (c *InitSymbols) Run(ctx context.Context) error {
	// scan symbols from API
	ss, err := c.securityAPI.ListSymbols(ctx, &api.SecurityFilter{
		Exchange: goutil.String(c.cfg.Exchange),
	})
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to list symbols from api, err: %v", err)
		return err
	}

	var (
		batchSize  = 1000
		batch      = make([]*entity.Security, 0, batchSize)
		retryCount = 10
		backOffMs  = 300
		count      = 0
	)

	for _, s := range ss {
		batch = append(batch, s)
		if len(batch) < batchSize {
			continue
		}

		if err := goutil.SyncRetry(ctx, func(ctx context.Context) error {
			return c.securityRepo.CreateMany(ctx, batch)
		}, retryCount, backOffMs); err != nil {
			log.Ctx(ctx).Error().Msgf("fail to save securities to mongo, err: %v", err)
			return err
		}

		count += len(batch)

		batch = make([]*entity.Security, 0, 1000)
	}

	if len(batch) > 0 {
		if err := goutil.SyncRetry(ctx, func(ctx context.Context) error {
			return c.securityRepo.CreateMany(ctx, batch)
		}, retryCount, backOffMs); err != nil {
			log.Ctx(ctx).Error().Msgf("fail to save securities to mongo, err: %v", err)
			return err
		}

		count += len(batch)
	}

	log.Ctx(ctx).Info().Msgf("inserted %v symbols, exchange: %v", count, c.cfg.Exchange)

	return nil
}

func (c *InitSymbols) Clean(ctx context.Context) error {
	return c.mongo.Close(ctx)
}
