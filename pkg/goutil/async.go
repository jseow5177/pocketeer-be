package goutil

import (
	"context"
	"errors"
	"sync"
	"time"
)

type noCancel struct {
	context.Context
}

func (c noCancel) Done() <-chan struct{} { return nil }

func (c noCancel) Deadline() (deadline time.Time, ok bool) { return }

func (c noCancel) Err() error { return nil }

// WithoutCancel returns a context that is never canceled.
func WithoutCancel(ctx context.Context) context.Context {
	return noCancel{ctx}
}

type composedContext struct {
	context.Context
	ctxs []context.Context
}

func (cc composedContext) Value(key interface{}) interface{} {
	for _, ctx := range cc.ctxs {
		if val := ctx.Value(key); val != nil {
			return val
		}
	}
	return cc.Context.Value(key)
}

func WithCancel(ctx context.Context, ctxs ...context.Context) (context.Context, context.CancelFunc) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	return composedContext{ctx, ctxs}, cancel
}

type Async struct {
	wg           sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
	initialDelay time.Duration
	retryNum     int
}

func NewAsync(initialDelay time.Duration, retryNum int) *Async {
	ctx, cancel := context.WithCancel(context.Background())
	m := &Async{
		initialDelay: initialDelay,
		retryNum:     retryNum,
		ctx:          ctx,
		cancel:       cancel,
	}
	return m
}

func (m *Async) wrap(cb func()) {
	m.wg.Add(1)
	go func() {
		cb()
		m.wg.Done()
	}()
}

func (m *Async) Wait() {
	m.wg.Wait()
}

func (m *Async) Stop(ctx context.Context) error {
	m.cancel()
	if waitTimeout(ctx, &m.wg) {
		return errors.New("context deadline exceeded")
	}
	return nil
}

func (m *Async) delay(ctx context.Context, delay time.Duration, fn func(context.Context)) {
	once := &sync.Once{}
	m.wrap(func() {
		c := make(chan time.Time, 1)
		t := time.AfterFunc(delay, func() {
			select {
			case <-ctx.Done():
				return
			default:
			}

			once.Do(func() {
				close(c)
				fn(ctx)
			})
		})
		select {
		case <-c:
		case <-ctx.Done():
			t.Stop()
			once.Do(func() {
				close(c)
				fn(ctx)
			}) // wait until fn returns
		}
	})
}

func (m *Async) Retry(ctx context.Context, fn func(context.Context) error) context.CancelFunc {
	var cancel context.CancelFunc
	ctx, cancel = WithCancel(m.ctx, ctx)
	m.retry(ctx, cancel, m.retryNum, m.initialDelay, fn)
	return cancel
}

func (m *Async) retry(
	ctx context.Context,
	cancel context.CancelFunc,
	retry int,
	interval time.Duration,
	fn func(context.Context) error,
) {
	select {
	case <-ctx.Done():
		return
	default:
		if retry <= 0 {
			cancel()
			return
		}
	}

	m.wrap(func() {
		if err := fn(ctx); err != nil {
			retry--
			if retry > 0 {
				m.delay(ctx, interval, func(context.Context) {
					interval *= 2
					m.retry(ctx, cancel, retry, interval, fn)
				})
			} else {
				cancel()
			}
		} else {
			cancel()
		}
	})
}

// waitTimeout waits for the WaitGroup for the specified max timeout.
// Returns true if timed out while waiting.
func waitTimeout(ctx context.Context, wg *sync.WaitGroup) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-ctx.Done():
		return true // timed out
	}
}
