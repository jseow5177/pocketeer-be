package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

type Job interface {
	Init(ctx context.Context) error
	Run(ctx context.Context) error
	Clean(ctx context.Context) error
}

var cmds = map[string]struct {
	desc string
	job  Job
}{}

func main() {
	ctx := context.Background()

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

	// Job init
	if err := cmd.job.Init(ctx); err != nil {
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
