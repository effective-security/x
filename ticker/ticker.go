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
	if t.ticker != nil {
		t.ticker.Stop()
		t.ticker = nil
	}
}

func (t *Ticker) GetStatus() string {
	return t.status
}

func (t *Ticker) SetStatus(status string) {
	t.status = status
}

func (t *Ticker) Count() int {
	return t.count
}
