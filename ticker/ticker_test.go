package ticker

import (
	"context"
	"sync/atomic"
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

	var runCount int64
	runFunc := func(ctx context.Context, status string, count int) {
		atomic.AddInt64(&runCount, 1)
	}

	ticker := New(ctx, 100*time.Millisecond, "running", runFunc)
	defer ticker.Stop()

	time.Sleep(350 * time.Millisecond)
	if atomic.LoadInt64(&runCount) < 3 {
		t.Errorf("expected runCount to be at least 3, got %d", atomic.LoadInt64(&runCount))
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
	finalCount := atomic.LoadInt64(&runCount)
	time.Sleep(150 * time.Millisecond)
	if atomic.LoadInt64(&runCount) != finalCount {
		t.Errorf("expected runCount to stop incrementing after cancel, got %d", atomic.LoadInt64(&runCount))
	}
}

func TestTickerWithStop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var runCount int64
	runFunc := func(ctx context.Context, status string, count int) {
		atomic.AddInt64(&runCount, 1)
	}

	ticker := New(ctx, 50*time.Millisecond, "running", runFunc)
	time.Sleep(350 * time.Millisecond)
	ticker.Stop()
	time.Sleep(50 * time.Millisecond)
	assert.Greater(t, atomic.LoadInt64(&runCount), int64(4))
}

// TestTickerCallbackNotCalledAfterCancel verifies that the run callback is not called
// when the context is cancelled, even if a tick occurs after cancellation.
func TestTickerCallbackNotCalledAfterCancel(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var runCount int64
	runFunc := func(ctx context.Context, status string, count int) {
		atomic.AddInt64(&runCount, 1)
	}

	// Create ticker with a very short interval to ensure we get multiple ticks
	ticker := New(ctx, 10*time.Millisecond, "running", runFunc)

	// Let it run for a bit to ensure it's working
	time.Sleep(50 * time.Millisecond)
	initialCount := atomic.LoadInt64(&runCount)
	assert.Greater(t, initialCount, int64(0), "Ticker should have executed at least once")

	// Cancel the context
	cancel()

	// Wait for a few more potential ticks
	time.Sleep(100 * time.Millisecond)

	// Verify that no more callbacks were executed after cancellation
	assert.Equal(t, initialCount, atomic.LoadInt64(&runCount), "Callback should not be called after context cancellation")

	// Verify the ticker is properly stopped
	ticker.Stop()
	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, initialCount, atomic.LoadInt64(&runCount), "Callback should still not be called after Stop()")
}
