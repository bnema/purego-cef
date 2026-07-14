//go:build cef_integration

package integration

import (
	"sync"
	"time"
)

const runtimeCloseTimeout = 10 * time.Second

// runtimeCleanup closes a created browser before shutting CEF down. Its hooks
// keep the failure path deterministic without requiring a provisioned runtime.
type runtimeCleanup struct {
	browserCreated    bool
	closeBrowser      func()
	onBeforeClose     <-chan struct{}
	doMessageLoopWork func()
	shutdown          func()
	closeTimeout      time.Duration
	now               func() time.Time
	sleep             func(time.Duration)
	once              sync.Once
}

func newRuntimeCleanup(shutdown, doMessageLoopWork func()) *runtimeCleanup {
	return &runtimeCleanup{
		shutdown:          shutdown,
		doMessageLoopWork: doMessageLoopWork,
		closeTimeout:      runtimeCloseTimeout,
		now:               time.Now,
		sleep:             time.Sleep,
	}
}

func (c *runtimeCleanup) closeAndShutdown() {
	c.once.Do(func() {
		now := c.now
		if now == nil {
			now = time.Now
		}
		sleep := c.sleep
		if sleep == nil {
			sleep = time.Sleep
		}
		if c.browserCreated && c.closeBrowser != nil {
			c.closeBrowser()
			deadline := now().Add(c.closeTimeout)
			for {
				select {
				case <-c.onBeforeClose:
					goto shutdown
				default:
				}
				if !now().Before(deadline) {
					goto shutdown
				}
				c.doMessageLoopWork()
				sleep(10 * time.Millisecond)
			}
		}
	shutdown:
		c.shutdown()
	})
}
