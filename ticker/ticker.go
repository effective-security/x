package ticker

import (
	"context"
	"time"
)

type Ticker struct {
	ticker *time.Ticker
	status string
	count  int
}

func New(ctx context.Context, every time.Duration, run func(ctx context.Context, status string, count int)) *Ticker {
	t := &Ticker{
		ticker: time.NewTicker(every),
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				t.Stop()
				return
			case <-t.ticker.C:
				t.count++
				run(ctx, t.status, t.count)
			}
		}
	}()

	return t
}

func (t *Ticker) Stop() {
	if t != nil && t.ticker != nil {
		t.ticker.Stop()
		//t.ticker = nil
	}
}

func (t *Ticker) GetStatus() string {
	if t == nil {
		return ""
	}
	return t.status
}

func (t *Ticker) SetStatus(status string) {
	if t != nil {
		t.status = status
	}
}

func (t *Ticker) Count() int {
	if t == nil {
		return 0
	}
	return t.count
}
