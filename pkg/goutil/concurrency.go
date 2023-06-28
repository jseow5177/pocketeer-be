package goutil

import (
	"context"

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
