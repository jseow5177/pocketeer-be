package syncquotes

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/dep/api"
	"github.com/jseow5177/pockteer-be/dep/api/finnhub"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/rs/zerolog/log"
)

const (
	maxApiCallsPerMin = 60
)

type JobConfig struct{}

type SyncQuotesCmd struct {
	mongo       *mongo.Mongo
	quoteRepo   repo.QuoteRepo
	holdingRepo repo.HoldingRepo
	securityAPI api.SecurityAPI
}

func (c *SyncQuotesCmd) initFlags() error {
	flagSet := flag.NewFlagSet(fmt.Sprintf("%s %s", filepath.Base(os.Args[0]), os.Args[1]), flag.ExitOnError)

	if err := flagSet.Parse(os.Args[2:]); err != nil {
		return err
	}

	return nil
}

func (c *SyncQuotesCmd) Init(ctx context.Context, cfg *config.Config) error {
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

	c.securityAPI = finnhub.NewFinnHubMgr(cfg.FinnHub)

	c.quoteRepo, err = mongo.NewQuoteMongo(ctx, c.mongo, c.securityAPI)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to init quote repo, err: %v", err)
		return err
	}

	c.holdingRepo = mongo.NewHoldingMongo(c.mongo)

	return nil
}

func (c *SyncQuotesCmd) Run(ctx context.Context) error {
	var (
		page  = 1
		limit = 1000
	)

	p := &repo.Paging{
		Limit: goutil.Uint32(uint32(limit)),
		Page:  goutil.Uint32(uint32(page)),
	}
	uniqueSymbols := make(map[string]struct{})

	for {
		hs, err := c.holdingRepo.GetMany(ctx, repo.NewHoldingFilter(
			repo.WithHoldingType(goutil.Uint32(uint32(entity.HoldingTypeDefault))),
			repo.WithHoldingPaging(p),
		))
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get holdings from repo, err: %v", err)
			return err
		}

		// deduplicate
		for _, h := range hs {
			uniqueSymbols[h.GetSymbol()] = struct{}{}
		}

		if len(hs) < limit {
			break
		}

		page++
		p.Page = goutil.Uint32(uint32(page))
	}

	var (
		quotes  = make([]*entity.Quote, 0)
		symbols = make([]string, 0)
	)
	for symbol := range uniqueSymbols {
		quote, err := c.securityAPI.GetLatestQuote(ctx, &api.SecurityFilter{
			Symbol: goutil.String(symbol),
		})
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get latest quote from api, err: %v", err)
			return err
		}
		quotes = append(quotes, quote)
		symbols = append(symbols, symbol)
	}

	for i, symbol := range symbols {
		if i > 0 && i%maxApiCallsPerMin == 0 {
			time.Sleep(time.Minute)
		}

		if err := c.quoteRepo.Upsert(ctx, repo.NewQuoteFilter(
			repo.WithQuoteSymbol(goutil.String(symbol)),
		), quotes[i]); err != nil {
			log.Ctx(ctx).Error().Msgf("fail to upsert quote to repo, symbol: %v, err: %v", symbol, err)
			return err
		}
	}

	log.Ctx(ctx).Info().Msgf("synced %v quotes, symbols: %v", len(quotes), symbols)

	return nil
}

func (c *SyncQuotesCmd) Clean(ctx context.Context) error {
	return c.mongo.Close(ctx)
}
