package ticker

import (
	"context"
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
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

	ticker.SetStatus("running")
	if ticker.GetStatus() != "running" {
		t.Errorf("expected status to be 'running', got %s", ticker.GetStatus())
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
