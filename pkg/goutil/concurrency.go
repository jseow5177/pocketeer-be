package goutil

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func ParallelizeWork(
	ctx context.Context,
	workCount int,
	maxThread int,
	workFn func(ctx context.Context, workNum int) error,
) error {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(maxThread)

	errIdx := make(chan int, workCount)
	for i := 0; i < workCount; i++ {
		i := i
		g.Go(func() error {
			if err := workFn(ctx, i); err != nil {
				errIdx <- i
				return err
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return fmt.Errorf("error encountered at worker %v, err: %v", <-errIdx, err)
	}

	return nil
}

func SyncRetry(ctx context.Context, fn func(context.Context) error, n, ms int) error {
	var err error
	for {
		if n <= 0 {
			break
		}
		err = fn(ctx)
		if err == nil {
			return nil
		}
		time.Sleep(time.Duration(ms) * time.Millisecond)
		n--
	}
	return err
}
