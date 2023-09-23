package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/pkg/logger"
	"github.com/rs/zerolog/log"

	ier "github.com/jseow5177/pockteer-be/cmd/job/init_exchange_rates"
	is "github.com/jseow5177/pockteer-be/cmd/job/init_symbols"
)

type Job interface {
	Init(ctx context.Context, cfg *config.Config) error
	Run(ctx context.Context) error
	Clean(ctx context.Context) error
}

var cmds = map[string]struct {
	desc string
	job  Job
}{
	"init_symbols": {
		desc: "scan symbols from third party API and save into mongo",
		job:  new(is.InitSymbols),
	},
	"init_exchange_rates": {
		desc: "get exchange rates from third party API and save into mongo",
		job:  new(ier.InitExchangeRates),
	},
}

func main() {
	opts := initOpts()

	ctx := logger.InitZeroLog(context.Background(), opts.LogLevel)

	flag.Usage = func() {
		usage := "usage: %s <command> [<args>]"
		for cmd, entry := range cmds {
			usage += fmt.Sprintf("\n\t%-15s\t%s", cmd, entry.desc)
		}
		usage += "\n"
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), usage, filepath.Base(os.Args[0]))
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(2)
	}

	cmd, ok := cmds[os.Args[1]]
	if !ok {
		flag.Usage()
		os.Exit(2)
	}

	cfg := config.NewConfig()
	if err := cfg.Subscribe(ctx, opts.ConfigFile); err != nil {
		log.Ctx(ctx).Fatal().Msgf("fail subscribe to config, err: %v", err)
		return
	}

	// Job init
	if err := cmd.job.Init(ctx, cfg); err != nil {
		log.Ctx(ctx).Fatal().Msgf("fail to init job, err: %v", err)
		return
	}

	// Job run
	if err := cmd.job.Run(ctx); err != nil {
		log.Ctx(ctx).Fatal().Msgf("fail to run job, err: %v", err)
		return
	}

	// Job stop
	if err := cmd.job.Clean(ctx); err != nil {
		log.Ctx(ctx).Fatal().Msgf("fail to clean job, err: %v", err)
		return
	}
}

func initOpts() *config.Option {
	opt := config.NewOptions()

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		opt.LogLevel = logLevel
	}

	if configFile := os.Getenv("CONFIG_FILE"); configFile != "" {
		opt.ConfigFile = configFile
	}

	return opt
}
