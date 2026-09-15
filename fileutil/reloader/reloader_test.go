package reloader

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testCheckInterval = 100 * time.Millisecond

type reloadEvent struct {
	filePath   string
	modifiedAt time.Time
}

func TestReloader(t *testing.T) {
	tick := make(chan time.Time)
	stopped := make(chan struct{})
	previousMakeTicker := makeTicker
	makeTicker = func(interval time.Duration) (func(), <-chan time.Time) {
		assert.Equal(t, testCheckInterval, interval)
		return func() { close(stopped) }, tick
	}
	t.Cleanup(func() { makeTicker = previousMakeTicker })

	file := filepath.Join(t.TempDir(), "reloaded.txt")
	require.NoError(t, os.WriteFile(file, []byte("initial"), 0o600))

	events := make(chan reloadEvent, 2)
	onChanged := func(filePath string, modifiedAt time.Time) {
		events <- reloadEvent{
			filePath:   filePath,
			modifiedAt: modifiedAt,
		}
	}

	startedAt := time.Now().UTC()
	r, err := NewReloader(file, testCheckInterval, onChanged)
	require.NoError(t, err)
	require.NotNil(t, r)

	require.NoError(t, r.Reload())
	first := receiveReloadEvent(t, events)
	assert.Equal(t, file, first.filePath)
	assert.True(t, first.modifiedAt.IsZero())
	assert.True(t, r.LoadedAt().After(startedAt))
	assert.Equal(t, uint32(1), r.LoadedCount())

	modifiedAt := time.Now().Add(time.Second)
	require.NoError(t, os.Chtimes(file, modifiedAt, modifiedAt))
	tick <- time.Now()
	second := receiveReloadEvent(t, events)
	assert.Equal(t, file, second.filePath)
	assert.True(t, modifiedAt.Equal(second.modifiedAt))
	assert.Equal(t, uint32(2), r.LoadedCount())

	require.NoError(t, r.Close())
	receiveSignal(t, stopped)
	err = r.Close()
	require.Error(t, err)
	assert.Equal(t, "already closed", err.Error())
}

func TestReloaderNilCallback(t *testing.T) {
	tick := make(chan time.Time)
	previousMakeTicker := makeTicker
	makeTicker = func(time.Duration) (func(), <-chan time.Time) {
		return func() {}, tick
	}
	t.Cleanup(func() { makeTicker = previousMakeTicker })

	r, err := NewReloader("unused", testCheckInterval, nil)
	require.NoError(t, err)
	require.NoError(t, r.Reload())
	assert.Equal(t, uint32(1), r.LoadedCount())
	require.NoError(t, r.Close())
}

func TestReloaderCloseConcurrent(t *testing.T) {
	tick := make(chan time.Time)
	previousMakeTicker := makeTicker
	makeTicker = func(time.Duration) (func(), <-chan time.Time) {
		return func() {}, tick
	}
	t.Cleanup(func() { makeTicker = previousMakeTicker })

	r, err := NewReloader("unused", testCheckInterval, nil)
	require.NoError(t, err)

	errs := make(chan error, 2)
	go func() { errs <- r.Close() }()
	go func() { errs <- r.Close() }()
	firstErr := <-errs
	secondErr := <-errs
	assert.True(t, (firstErr == nil) != (secondErr == nil))
	if firstErr != nil {
		assert.Equal(t, "already closed", firstErr.Error())
	}
	if secondErr != nil {
		assert.Equal(t, "already closed", secondErr.Error())
	}
}

func TestReloaderCloseNil(t *testing.T) {
	var r *Reloader
	assert.NoError(t, r.Close())
}

func receiveReloadEvent(t *testing.T, events <-chan reloadEvent) reloadEvent {
	t.Helper()
	select {
	case event := <-events:
		return event
	case <-time.After(time.Second):
		require.FailNow(t, "timed out waiting for reload callback")
		return reloadEvent{}
	}
}

func receiveSignal(t *testing.T, signal <-chan struct{}) {
	t.Helper()
	select {
	case <-signal:
	case <-time.After(time.Second):
		require.FailNow(t, "timed out waiting for signal")
	}
}
