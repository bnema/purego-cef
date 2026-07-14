//go:build cef_integration

package integration

import (
	"reflect"
	"testing"
	"time"
)

func TestRuntimeCleanupClosesBrowserBeforeSingleShutdown(t *testing.T) {
	closed := make(chan struct{}, 1)
	var events []string
	cleanup := runtimeCleanup{
		browserCreated: true,
		closeBrowser:   func() { events = append(events, "close") },
		onBeforeClose:  closed,
		doMessageLoopWork: func() {
			events = append(events, "pump")
			closed <- struct{}{}
		},
		shutdown:     func() { events = append(events, "shutdown") },
		closeTimeout: time.Second,
	}

	// This is the deferred failure path after OnAfterCreated has fired.
	cleanup.closeAndShutdown()
	cleanup.closeAndShutdown()

	if want := []string{"close", "pump", "shutdown"}; !reflect.DeepEqual(events, want) {
		t.Fatalf("cleanup events = %v, want %v", events, want)
	}
}

func TestRuntimeCleanupBoundsCloseWaitBeforeSingleShutdown(t *testing.T) {
	var (
		now        = time.Unix(0, 0)
		closeCalls int
		pumpCalls  int
		shutdowns  int
	)
	cleanup := runtimeCleanup{
		browserCreated:    true,
		closeBrowser:      func() { closeCalls++ },
		onBeforeClose:     make(chan struct{}),
		doMessageLoopWork: func() { pumpCalls++ },
		shutdown:          func() { shutdowns++ },
		closeTimeout:      20 * time.Millisecond,
		now:               func() time.Time { return now },
		sleep:             func(d time.Duration) { now = now.Add(d) },
	}

	cleanup.closeAndShutdown()

	if closeCalls != 1 || pumpCalls != 2 || shutdowns != 1 {
		t.Fatalf("close=%d pumps=%d shutdowns=%d, want close=1 pumps=2 shutdowns=1", closeCalls, pumpCalls, shutdowns)
	}
}
