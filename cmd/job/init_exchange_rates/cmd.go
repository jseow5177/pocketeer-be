package initexchangerates

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/dep/api"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"

	exchangeratehost "github.com/jseow5177/pockteer-be/dep/api/exchange_rate_host"
)

type JobConfig struct {
	Date string
	From string
	To   string
	Unit string
}

type InitExchangeRates struct {
	cfg JobConfig

	mongo *mongo.Mongo

	exchangeRateRepo repo.ExchangeRateRepo
	exchangeRateAPI  api.ExchangeRateAPI

	fromCurrencies []string
	toCurrencies   []string
}

func (c *InitExchangeRates) initFlags() error {
	flagSet := flag.NewFlagSet(fmt.Sprintf("%s %s", filepath.Base(os.Args[0]), os.Args[1]), flag.ExitOnError)

	flagSet.StringVar(&c.cfg.Date, "date", util.FormatDate(time.Now()), "date of exchange rate, format: 20220202")
	flagSet.StringVar(&c.cfg.From, "from", "", "comma-separated currencies, eg: SGD,MYR")
	flagSet.StringVar(&c.cfg.To, "to", "", "comma-separated currencies, eg: SGD,MYR")

	if err := flagSet.Parse(os.Args[2:]); err != nil {
		return err
	}

	// default currencies to use
	currencies := make([]string, 0)
	for currency := range entity.Currencies {
		currencies = append(currencies, currency)
	}
	c.fromCurrencies = currencies
	c.toCurrencies = currencies

	if c.cfg.From != "" {
		fromCurrencies := strings.Split(c.cfg.From, ",")
		for _, currency := range fromCurrencies {
			if err := entity.CheckCurrency(currency); err != nil {
				return err
			}
		}
		c.fromCurrencies = fromCurrencies
	}

	if c.cfg.To != "" {
		toCurrencies := strings.Split(c.cfg.To, ",")
		for _, currency := range toCurrencies {
			if err := entity.CheckCurrency(currency); err != nil {
				return err
			}
		}
		c.toCurrencies = toCurrencies
	}

	return nil
}

func (c *InitExchangeRates) Init(ctx context.Context, cfg *config.Config) error {
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

	c.exchangeRateRepo, err = mongo.NewExchangeRateMongo(ctx, c.mongo)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to init exchange rate repo, err: %v", err)
		return err
	}

	c.exchangeRateAPI = exchangeratehost.NewExchangeRateHostMgr(cfg.ExchangeRateHost)

	return nil
}

func (c *InitExchangeRates) Run(ctx context.Context) error {
	date, err := util.ParseDate(c.cfg.Date)
	if err != nil {
		return fmt.Errorf("fail to parse date, date: %v, err: %v", c.cfg.Date, err)
	}

	exchangeRates := make([]*entity.ExchangeRate, 0)
	for _, fromCurrency := range c.fromCurrencies {
		symbols := make([]string, 0)
		for _, toCurrency := range c.toCurrencies {
			if toCurrency != fromCurrency { // redundant to store 1:1
				symbols = append(symbols, toCurrency)
			}
		}

		if len(symbols) == 0 {
			continue
		}

		erf := api.NewExchangeRateFilter(
			util.FormatDate(date),
			api.WithExchangeRateBase(goutil.String(fromCurrency)),
			api.WithExchangeRateCurrencies(symbols...),
		)

		ers, err := c.exchangeRateAPI.GetExchangeRates(ctx, erf)
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get exchange rates, err: %v", err)
			return err
		}

		exchangeRates = append(exchangeRates, ers...)
	}

	if len(exchangeRates) == 0 {
		log.Ctx(ctx).Info().Msg("no exchange rates created")
		return nil
	}

	ids, err := c.exchangeRateRepo.CreateMany(ctx, exchangeRates)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to create exchange rates in repo, err: %v", err)
		return err
	}

	log.Ctx(ctx).Info().Msgf("inserted %v exchange rates", len(ids))

	return nil
}

func (c *InitExchangeRates) Clean(ctx context.Context) error {
	return c.mongo.Close(ctx)
}
