package ticker

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNilTicker verifies nil ticker methods do not panic and return defaults.
func TestNilTicker(t *testing.T) {
	t.Parallel()
	var tkr *Ticker[string]
	assert.Empty(t, tkr.GetValue())
	assert.Equal(t, 0, tkr.Count())
	// Setting status on nil should be a no-op and not panic
	tkr.SetValue("test")
	assert.Empty(t, tkr.GetValue())
}

func TestTickerWithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runCount := 0
	runFunc := func(ctx context.Context, status string, count int) {
		runCount++
	}

	ticker := New(ctx, 100*time.Millisecond, runFunc)
	defer ticker.Stop()

	time.Sleep(350 * time.Millisecond)
	if runCount < 3 {
		t.Errorf("expected runCount to be at least 3, got %d", runCount)
	}

	ticker.SetValue("running")
	if ticker.GetValue() != "running" {
		t.Errorf("expected status to be 'running', got %s", ticker.GetValue())
	}

	if ticker.Count() < 3 {
		t.Errorf("expected count to be at least 3, got %d", ticker.Count())
	}

	cancel()
	time.Sleep(150 * time.Millisecond)
	finalCount := runCount
	time.Sleep(150 * time.Millisecond)
	if runCount != finalCount {
		t.Errorf("expected runCount to stop incrementing after cancel, got %d", runCount)
	}
}

func TestTickerWithStop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runCount := 0
	runFunc := func(ctx context.Context, status string, count int) {
		runCount++
	}

	ticker := New(ctx, 50*time.Millisecond, runFunc)
	time.Sleep(350 * time.Millisecond)
	ticker.Stop()
	time.Sleep(50 * time.Millisecond)
	assert.Greater(t, runCount, 0)
}
