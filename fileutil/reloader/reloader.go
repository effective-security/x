package reloader

import (
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/effective-security/xlog"
)

var logger = xlog.NewPackageLogger("github.com/effective-security/x/fileutil", "reloader")

// Wrap time.Tick so we can override it in tests.
var makeTicker = func(interval time.Duration) (func(), <-chan time.Time) {
	t := time.NewTicker(interval)
	return t.Stop, t.C
}

// OnChangedFunc is a called when the file has been modified
type OnChangedFunc func(filePath string, modifiedAt time.Time)

// Reloader keeps necessary info to provide reloaded certificate
type Reloader struct {
	lock           sync.RWMutex
	loadedAt       time.Time
	count          uint32
	filePath       string
	fileModifiedAt time.Time
	onChangedFunc  OnChangedFunc
	inProgress     bool
	stopChan       chan struct{}
	closed         bool
}

// NewReloader return an instance of the file re-loader
func NewReloader(filePath string, checkInterval time.Duration, onChangedFunc OnChangedFunc) (*Reloader, error) {
	stopChan := make(chan struct{})
	result := &Reloader{
		filePath:      filePath,
		onChangedFunc: onChangedFunc,
		stopChan:      stopChan,
	}

	logger.KV(xlog.INFO, "status", "started", "file", filePath)

	tickerStop, tickChan := makeTicker(checkInterval)
	go func() {
		for {
			select {
			case <-stopChan:
				tickerStop()
				logger.KV(xlog.INFO, "status", "closed", "count", result.LoadedCount(), "file", filePath)
				return
			case <-tickChan:
				fi, err := os.Stat(filePath)
				if err != nil {
					logger.KV(xlog.WARNING, "reason", "stat", "file", filePath, "err", err)
					continue
				}
				modTime := fi.ModTime()
				result.lock.Lock()
				modified := modTime.After(result.fileModifiedAt)
				if modified {
					result.fileModifiedAt = modTime
				}
				result.lock.Unlock()
				if modified {
					if err := result.Reload(); err != nil {
						logger.KV(xlog.ERROR, "err", err)
					}
				}
			}
		}
	}()
	return result, nil
}

// Reload will explicitly call the callback function
func (k *Reloader) Reload() error {
	k.lock.Lock()
	if k.inProgress {
		k.lock.Unlock()
		return nil
	}

	k.inProgress = true
	modifiedAt := k.fileModifiedAt
	path := k.filePath
	cb := k.onChangedFunc
	atomic.AddUint32(&k.count, 1)
	k.loadedAt = time.Now().UTC()
	k.inProgress = false
	k.lock.Unlock()

	if cb != nil {
		go cb(path, modifiedAt)
	}

	return nil
}

// LoadedAt return the last time when the pair was loaded
func (k *Reloader) LoadedAt() time.Time {
	k.lock.RLock()
	defer k.lock.RUnlock()

	return k.loadedAt
}

// LoadedCount returns the number of times the pair was loaded from disk
func (k *Reloader) LoadedCount() uint32 {
	return atomic.LoadUint32(&k.count)
}

// Close will close the reloader and release its resources
func (k *Reloader) Close() error {
	if k == nil {
		return nil
	}

	k.lock.Lock()
	defer k.lock.Unlock()

	if k.closed {
		return errors.New("already closed")
	}

	k.closed = true
	close(k.stopChan)

	return nil
}
