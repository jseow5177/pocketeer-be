package goutil

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"
)

func ParallelizeWork(
	ctx context.Context,
	workCount int,
	maxThread int,
	workFn func(ctx context.Context, workNum int) error,
) error {
	var (
		g      = new(errgroup.Group)
		wgChan = make(chan struct{}, maxThread)
	)

	for i := 0; i < workCount; i++ {
		select {
		case wgChan <- struct{}{}:
		case <-ctx.Done():
			return ctx.Err()
		}

		i := i
		g.Go(func() error {
			defer func() {
				<-wgChan
			}()

			return workFn(ctx, i)
		})
	}

	return g.Wait()
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
