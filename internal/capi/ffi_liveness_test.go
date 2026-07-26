package capi

import (
	"runtime"
	"sync"
	"testing"
	"time"
	"unsafe"

	"github.com/bnema/purego"
)

const (
	loadedCISynchronizationTimeout = 10 * time.Second
	finalizerPollingInterval       = 10 * time.Millisecond
)

type ffiTrackedArgument struct {
	marker uintptr
}

type ffiCallbackObservation struct {
	self       uintptr
	event      uintptr
	mouseleave uintptr
}

func TestCallSendMouseMoveEventLivenessABI(t *testing.T) {
	testCases := []struct {
		name       string
		mouseleave uintptr
		withEvent  bool
	}{
		{name: "zero", mouseleave: 0},
		{name: "pointer", mouseleave: 1, withEvent: true},
		{name: "max", mouseleave: ^uintptr(0), withEvent: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testBlockingFFICall(t, testCase.withEvent, testCase.mouseleave)
		})
	}
}

func testBlockingFFICall(t *testing.T, withEvent bool, mouseleave uintptr) {
	t.Helper()

	entered := make(chan struct{})
	release := make(chan struct{})
	observed := make(chan ffiCallbackObservation, 1)
	callDone := make(chan struct{})
	finalized := make(chan struct{}, 1)
	finalizedBeforeCallbackRelease := make(chan struct{}, 1)
	var releaseOnce sync.Once
	releaseCallback := func() { releaseOnce.Do(func() { close(release) }) }

	callback := purego.NewCallback(func(self, event, callbackMouseleave uintptr) {
		close(entered)
		select {
		case <-release:
		case <-time.After(loadedCISynchronizationTimeout):
			return
		}
		observed <- ffiCallbackObservation{
			self:       self,
			event:      event,
			mouseleave: callbackMouseleave,
		}
	})
	t.Cleanup(func() {
		if err := purego.UnrefCallback(callback); err != nil {
			t.Errorf("unref callback: %v", err)
		}
	})
	// Registered after callback cleanup so LIFO cleanup releases a callback
	// before its trampoline is unreferenced when a test assertion fails.
	t.Cleanup(releaseCallback)

	host := &CEFBrowserHostT{}
	host.OverrideSendMouseMoveEvent(callback)
	expectedSelf := uintptr(unsafe.Pointer(host))

	var argument *ffiTrackedArgument
	var event unsafe.Pointer
	if withEvent {
		argument = &ffiTrackedArgument{marker: 0xc0ffee}
		runtime.SetFinalizer(argument, func(*ffiTrackedArgument) {
			select {
			case <-release:
				finalized <- struct{}{}
			default:
				finalizedBeforeCallbackRelease <- struct{}{}
			}
		})
		event = unsafe.Pointer(argument)
	}
	expectedEvent := uintptr(event)

	go func(arg unsafe.Pointer) {
		host.CallSendMouseMoveEvent(arg, mouseleave)
		close(callDone)
	}(event)
	argument = nil
	event = nil

	if !waitSignal(entered, loadedCISynchronizationTimeout) {
		releaseCallback()
		t.Fatal("timed out waiting for FFI callback entry")
	}

	forceFinalizerOpportunity(t)
	if withEvent {
		select {
		case <-finalizedBeforeCallbackRelease:
			releaseCallback()
			t.Fatal("argument finalizer ran while the FFI callback was blocked")
		default:
		}
	}

	releaseCallback()
	var got ffiCallbackObservation
	select {
	case got = <-observed:
	case <-time.After(loadedCISynchronizationTimeout):
		t.Fatal("timed out waiting for callback argument observation")
	}
	if got.self != expectedSelf || got.event != expectedEvent || got.mouseleave != mouseleave {
		t.Errorf("callback arguments = {%#x, %#x, %#x}, want {%#x, %#x, %#x}",
			got.self, got.event, got.mouseleave, expectedSelf, expectedEvent, mouseleave)
	}
	if !waitSignal(callDone, loadedCISynchronizationTimeout) {
		t.Fatal("timed out waiting for FFI call return")
	}

	if withEvent {
		waitForFinalizer(t, finalized, finalizedBeforeCallbackRelease)
	}
}

func forceFinalizerOpportunity(t *testing.T) {
	t.Helper()

	pressureDone := make(chan struct{})
	go func() {
		defer close(pressureDone)
		for range 8 {
			pressure := make([]byte, 1<<20)
			pressure[0] = 1
			runtime.GC()
			runtime.Gosched()
		}
	}()
	if !waitSignal(pressureDone, loadedCISynchronizationTimeout) {
		t.Fatal("timed out applying GC and allocation pressure")
	}
}

func waitForFinalizer(t *testing.T, finalized, finalizedBeforeCallbackRelease <-chan struct{}) {
	t.Helper()

	deadline := time.NewTimer(loadedCISynchronizationTimeout)
	defer deadline.Stop()
	poll := time.NewTicker(finalizerPollingInterval)
	defer poll.Stop()

	runtime.GC()
	for {
		select {
		case <-finalizedBeforeCallbackRelease:
			t.Fatal("argument finalizer ran while the FFI callback was blocked")
		case <-finalized:
			return
		case <-deadline.C:
			t.Fatal("timed out waiting for argument finalizer after FFI call return")
		case <-poll.C:
			runtime.GC()
			runtime.Gosched()
		}
	}
}

func waitSignal(signal <-chan struct{}, timeout time.Duration) bool {
	select {
	case <-signal:
		return true
	case <-time.After(timeout):
		return false
	}
}
