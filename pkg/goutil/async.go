package goutil

import (
	"context"
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
