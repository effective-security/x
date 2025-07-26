package ticker

import (
	"context"
	"sync"
	"time"
)

type Ticker[T any] struct {
	ticker *time.Ticker
	value  T
	count  int
	lock   sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

func New[T any](ctx context.Context, every time.Duration, value T, run func(ctx context.Context, value T, count int)) *Ticker[T] {
	tickerCtx, cancel := context.WithCancel(ctx)
	t := &Ticker[T]{
		ticker: time.NewTicker(every),
		ctx:    tickerCtx,
		value:  value,
		cancel: cancel,
	}

	go func() {
		for {
			select {
			case <-tickerCtx.Done():
				t.close()
				return
			case <-t.ticker.C:
				// Check if context is cancelled before executing the callback
				if tickerCtx.Err() != nil {
					return
				}

				t.lock.Lock()
				t.count++
				value := t.value
				count := t.count
				t.lock.Unlock()

				run(tickerCtx, value, count)
			}
		}
	}()

	return t
}

func (t *Ticker[T]) Stop() {
	if t == nil {
		return
	}

	// Cancel the context first to stop the goroutine
	if t.cancel != nil {
		t.cancel()
	}
}

func (t *Ticker[T]) close() {
	if t == nil {
		return
	}
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.ticker != nil {
		t.ticker.Stop()
		t.ticker = nil
	}
}

func (t *Ticker[T]) GetValue() T {
	if t == nil {
		var empty T
		return empty
	}
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.value
}

func (t *Ticker[T]) SetValue(value T) {
	if t == nil {
		return
	}
	t.lock.Lock()
	defer t.lock.Unlock()
	t.value = value
}

func (t *Ticker[T]) Count() int {
	if t == nil {
		return 0
	}
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.count
}
